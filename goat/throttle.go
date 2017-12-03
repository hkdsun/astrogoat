package goat

import "database/sql"

type Throttler interface {
	Call()
}

type NaiveThrottler struct {
	DB *sql.DB
}

func (t *NaiveThrottler) Call() {
	t.DB.Query("SELECT * FROM meta.heartbeat")
}
