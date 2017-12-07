package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hkdsun/astrogoat/goat"
)

func mysqlConfig(user, password, addr string) *mysql.Config {
	return &mysql.Config{
		User:   user,
		Passwd: password,

		Net:  "tcp",
		Addr: addr,
	}
}

func insertLoad(db *sql.DB) (string, []interface{}) {
	randId := rand.Intn(9999999999)
	return "INSERT IGNORE INTO test.test VALUES(NULL,?)", []interface{}{randId}
}

func createTestDbAndTable(db *sql.DB) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS test")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS test")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test.test (id INT AUTO_INCREMENT, data INT, PRIMARY KEY(id))")
	if err != nil {
		return err
	}

	return nil
}

func naiveThrottle() *goat.NaiveThrottle {
	slaveConfig := mysqlConfig("root", "", "127.0.0.1:22002")
	slaveDb, err := sql.Open("mysql", slaveConfig.FormatDSN())
	if err != nil {
		panic("Failed to connect to slaveDb")
	}

	naiveThrottle := &goat.NaiveThrottle{
		DB: slaveDb,
	}

	go naiveThrottle.Run()

	return naiveThrottle
}

func pidThrottle() *goat.PidThrottle {
	slaveConfig := mysqlConfig("root", "", "127.0.0.1:22002")
	slaveDb, err := sql.Open("mysql", slaveConfig.FormatDSN())
	if err != nil {
		panic("Failed to connect to slaveDb")
	}

	pidThrottle := goat.NewPidThrottle(slaveDb)

	go pidThrottle.Run()

	return pidThrottle
}

func slowStart() *goat.SlowStartThrottle {
	slaveConfig := mysqlConfig("root", "", "127.0.0.1:22002")
	slaveDb, err := sql.Open("mysql", slaveConfig.FormatDSN())
	if err != nil {
		panic("Failed to connect to slaveDb")
	}

	slowStartThrottle := &goat.SlowStartThrottle{
		DB:               slaveDb,
		CurrentSleep:     0 * time.Second,
		BestSleep:        10 * time.Second,
		IncreaseStepSize: 200 * time.Millisecond,
		DecreaseStepSize: 50 * time.Millisecond,
		LagThreshold:     500 * time.Millisecond,
		CacheDuration:    2000 * time.Millisecond,
	}

	return slowStartThrottle
}

func main() {
	dbConfig := mysqlConfig("root", "", "127.0.0.1:21001")
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic("Failed to connect to db")
	}

	generators := []goat.LoadGenerator{
		&goat.BatchGenerator{
			BatchSize: 50,
			QueryFunc: insertLoad,
		},
	}

	job := &goat.Job{
		JobConfig: &goat.JobConfig{
			SetupFunc:      createTestDbAndTable,
			Routines:       9,
			LoadGenerators: generators,
			Duration:       120 * time.Second,
			Interval:       50 * time.Millisecond,
			DB:             db,
			Throttler:      slowStart(),
		},
	}

	err = job.Run()
	if err != nil {
		panic(fmt.Sprintf("Job failed: %+v", err))
	}

}
