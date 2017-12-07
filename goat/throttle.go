package goat

import (
	"database/sql"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	cacheDuration      = 5000 * time.Millisecond
	unhealthyThreshold = 3000000
)

type Throttler interface {
	Call()
}

type NaiveThrottle struct {
	DB *sql.DB

	curLag        int
	mut           *sync.RWMutex
	lastRead      time.Time
	lastThrottled time.Time
}

func (t *NaiveThrottle) updateLag() {
	var lag int
	row := t.DB.QueryRow("SELECT MAX(TIMESTAMPDIFF(MICROSECOND, ts, NOW())) FROM meta.heartbeat")

	err := row.Scan(&lag)
	if err != nil {
		logrus.WithError(err).Error("Could not get lag")
	}

	t.setCurLag(lag)
	t.lastRead = time.Now()

	logrus.WithField("lag", t.getCurLag()).Info("updated lag")
}

func (t *NaiveThrottle) getCurLag() int {
	t.mut.RLock()
	defer t.mut.RUnlock()
	return t.curLag
}

func (t *NaiveThrottle) setCurLag(lag int) {
	t.mut.Lock()
	defer t.mut.Unlock()
	t.curLag = lag
}

func (t *NaiveThrottle) Run() {
	t.mut = &sync.RWMutex{}

	for {
		if time.Now().Sub(t.lastRead) >= cacheDuration {
			logrus.Info("updating lag")
			t.updateLag()
		}
	}
}

func (t *NaiveThrottle) Call() {
	if t.getCurLag() > unhealthyThreshold {
		t.awaitHealthy()
	}
}

func (t *NaiveThrottle) awaitHealthy() {
	t.lastThrottled = time.Now()
	for {
		if t.getCurLag() > unhealthyThreshold {
			logrus.Debug("Sleeping since lag is %d", t.getCurLag())
			time.Sleep(cacheDuration + time.Duration(rand.Int31n(5000))*time.Millisecond)
		} else {
			return
		}
	}
}
