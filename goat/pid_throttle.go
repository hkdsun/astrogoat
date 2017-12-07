package goat

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

type PidThrottle struct {
	DB    *sql.DB
	pid   *PIDController
	sleep time.Duration
}

func NewPidThrottle(db *sql.DB) *PidThrottle {
	control := NewPIDController(-3.1, 0, 0)
	control.SetOutputLimits(0.0, 30000000.0)
	control.Set(1000000)

	throttle := &PidThrottle{
		DB:  db,
		pid: control,
	}

	return throttle
}

func (t *PidThrottle) update() {
	var lag float64
	row := t.DB.QueryRow("SELECT MAX(TIMESTAMPDIFF(MICROSECOND, ts, NOW())) FROM meta.heartbeat")

	err := row.Scan(&lag)
	if err != nil {
		logrus.WithError(err).Error("Could not get lag")
	}

	if lag < 0 {
		lag = 0
	}

	t.sleep = time.Duration(t.pid.Update(lag)) * time.Microsecond
	// logrus.WithField("lag", lag).WithField("sleep", t.sleep).Info("updated lag")
}

func (t *PidThrottle) Run() {
	for {
		t.update()
		time.Sleep(1 * time.Second)
	}
}

func (t *PidThrottle) Call() {
	// sleep(t.sleep * [0,1])
	// randSleep := t.sleep.Seconds() * rand.Float64()
	// s := time.Duration(randSleep) * time.Second

	s := time.Duration(rand.Float64()*1000) * time.Millisecond
	logrus.WithField("rand_sleep", s).WithField("sleep", t.sleep).Info("sleeping")

	time.Sleep(s)

	time.Sleep(t.sleep)
}
