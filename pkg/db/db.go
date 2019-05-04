package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/ryantking/marina/pkg/config"
	"upper.io/db.v3/lib/sqlbuilder"
)

const (
	mysqlParams = "parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
)

var (
	db  sqlbuilder.Database
	dbL sync.Mutex
)

// Get returns the connection to the database
func Get() (sqlbuilder.Database, error) {
	dbL.Lock()
	defer dbL.Unlock()
	if db != nil {
		return db, nil
	}

	cfg := config.Get()
	connStr := cfg.DB.DSN
	if cfg.DB.Type == "mysql" {
		connStr = fmt.Sprintf("%s?%s", connStr, mysqlParams)
	}
	sqlDB, err := sql.Open(cfg.DB.Type, connStr)
	if err != nil {
		return nil, err
	}
	sess, err := sqlbuilder.New(cfg.DB.Type, sqlDB)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(sess.Context(), cfg.DB.Timeout)
	defer cancel()

	err = sess.WithContext(ctx).Ping()
	if err != nil {
		return nil, err
	}

	db = sess
	return db, nil
}

// Close closes the connection to the database
func Close() error {
	dbL.Lock()
	defer dbL.Unlock()
	if db == nil {
		return nil
	}

	err := db.Close()
	if err != nil {
		return err
	}

	db = nil
	return nil
}