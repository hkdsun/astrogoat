package goat

import (
	"context"
	"database/sql"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type JobConfig struct {
	Routines       int
	LoadGenerators []LoadGenerator
	Duration       time.Duration
	Interval       time.Duration
	DB             *sql.DB
	SetupFunc      func(*sql.DB) error
	Throttler      Throttler
}

type Job struct {
	*JobConfig
}

func (j *Job) Run() error {
	if j.SetupFunc != nil {
		err := j.SetupFunc(j.DB)
		if err != nil {
			return err
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(j.Routines)

	for i := 0; i < j.Routines; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), j.Duration)
			defer cancel()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					j.runRandomLoadGenerator()
					time.Sleep(j.Interval)
				}
			}
		}()
	}
	wg.Wait()

	return nil
}

func (j *Job) runRandomLoadGenerator() {
	gen := rand.Intn(len(j.LoadGenerators))
	err := j.LoadGenerators[gen].Apply(j.DB)
	if err != nil {
		logrus.WithError(err).Error("could not apply load")
	}
}
