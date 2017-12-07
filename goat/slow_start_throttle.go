package goat

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	S_STARTING int = iota
	S_TUNING
)

var (
	minSleep = 100 * time.Millisecond
)

type SlowStartThrottle struct {
	DB               *sql.DB
	CurrentSleep     time.Duration
	BestSleep        time.Duration
	IncreaseStepSize time.Duration
	DecreaseStepSize time.Duration
	LagThreshold     time.Duration
	CacheDuration    time.Duration

	lagAvg      time.Duration
	lastUpdated time.Time
	state       int
}

func (t *SlowStartThrottle) updateLagAvg(newSample time.Duration) {
	N := 50 * time.Second / t.CacheDuration
	fmt.Printf("N = %+v\n", N)
	t.lagAvg -= t.lagAvg / N
	t.lagAvg += newSample / N
}

func (t *SlowStartThrottle) getLag() time.Duration {
	var lagi int
	row := t.DB.QueryRow("SELECT MAX(TIMESTAMPDIFF(MICROSECOND, ts, NOW())) FROM meta.heartbeat")

	err := row.Scan(&lagi)
	if err != nil {
		logrus.WithError(err).Error("Could not get lagi")
	}

	lag := time.Duration(lagi) * time.Microsecond

	if lag <= 0 {
		return 0
	}

	return lag
}

func (t *SlowStartThrottle) unhealthy() bool {
	return t.getLag() > t.LagThreshold
}

func (t *SlowStartThrottle) Call() {
	time.Sleep(t.CurrentSleep)

	if time.Since(t.lastUpdated) <= t.CacheDuration {
		return
	}

	lag := t.getLag()
	t.updateLagAvg(lag)
	fmt.Printf("t.lagAvg = %+v\n", t.lagAvg)
	fmt.Printf("t.lag = %+v\n", lag)
	fmt.Printf("t.CurrentSleep = %+v\n", t.CurrentSleep)

	if t.unhealthy() {
		t.CurrentSleep += lag/8 + t.lagAvg*2
	} else {
		lagBudget := t.LagThreshold - lag
		if lagBudget > 0 {
			t.CurrentSleep -= lagBudget
		}
		t.CurrentSleep -= t.DecreaseStepSize
	}

	if t.CurrentSleep < minSleep {
		t.CurrentSleep = minSleep
	}
	t.lastUpdated = time.Now()
}
