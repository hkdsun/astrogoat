package goat

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type LoadGenerator interface {
	Apply(db *sql.DB) error
}

type SimpleLoadGenerator struct {
	QueryFunc func(db *sql.DB) (string, []interface{})
	Query     string
}

func (s *SimpleLoadGenerator) Apply(db *sql.DB) error {
	args := []interface{}{}
	query := s.Query

	if s.QueryFunc != nil {
		query, args = s.QueryFunc(db)
	}

	logrus.
		WithField("query", query).
		WithField("args", args).
		Debug("Executing")

	_, err := db.Exec(query, args...)
	return err
}
