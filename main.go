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
	randId := rand.Intn(99999)
	return "INSERT IGNORE INTO test.test VALUES(?)", []interface{}{randId}
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test.test (id INT, PRIMARY KEY(id))")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	dbConfig := mysqlConfig("root", "", "127.0.0.1:21001")
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic("Failed to connect to db")
	}

	naiveThrottle := goat.NaiveThrottle{
		DB: db,
	}

	generators := []goat.LoadGenerator{
		&goat.SimpleLoadGenerator{
			QueryFunc: insertLoad,
		},
	}

	job := &goat.Job{
		JobConfig: &goat.JobConfig{
			SetupFunc:      createTestDbAndTable,
			Routines:       100,
			LoadGenerators: generators,
			Duration:       30 * time.Second,
			Interval:       20 * time.Millisecond,
			DB:             db,
			Throttler:      naiveThrottle,
		},
	}

	err = job.Run()
	if err != nil {
		panic(fmt.Sprintf("Job failed: %+v", err))
	}

}
