package goat

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type BatchGenerator struct {
	BatchSize int
	QueryFunc func(db *sql.DB) (string, []interface{})
	Query     string
}

func (s *BatchGenerator) SetBatchSize(size int) {
	s.BatchSize = size
}

func (s *BatchGenerator) Apply(db *sql.DB) error {
	args := []interface{}{}
	query := s.Query

	if s.QueryFunc != nil {
		query, args = s.QueryFunc(db)
	}

	logrus.
		WithField("query", query).
		WithField("args", args).
		Debug("Executing")

	for batch := 0; batch < s.BatchSize; batch++ {
		_, err := db.Exec(query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
