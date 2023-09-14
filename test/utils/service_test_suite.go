package test_utils

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type Suite struct {
	sqlDb *sql.DB
	DB    *gorm.DB
	Mock  sqlmock.Sqlmock
}

func NewSuite(t *testing.T) *Suite {
	s := &Suite{}
	var err error
	s.sqlDb, s.Mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	if s.sqlDb == nil {
		t.Error("mock db is null")
	}
	if s.Mock == nil {
		t.Error("sqlmock is null")
	}
	if s.DB, err = gorm.Open(
		postgres.New(postgres.Config{Conn: s.sqlDb, DriverName: "postgres"}),
		&gorm.Config{},
	); err != nil {
		panic(err) // Error here
	}
	return s
}

func (s *Suite) Close() {
	_ = s.sqlDb.Close()
}
