package db

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ryantking/marina/pkg/config"
)

const (
	mysqlParams = "parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
)

var (
	db  *gorm.DB
	dbL sync.Mutex
)

// Get returns the connection to the database
func Get() *gorm.DB {
	dbL.Lock()
	defer dbL.Unlock()
	if db != nil {
		return db
	}

	cfg := config.Get()
	connStr := cfg.DB.DSN
	if cfg.DB.Type == "mysql" {
		connStr = fmt.Sprintf("%s?%s", connStr, mysqlParams)
	}
	d, err := gorm.Open(cfg.DB.Type, connStr)
	if err != nil {
		panic(fmt.Sprintf("failed to get database connection: %s", err.Error()))
	}
	db = d
	return db
}

// Close closes the connection to the database
func Close() {
	dbL.Lock()
	defer dbL.Unlock()
	if db == nil {
		return
	}

	err := db.Close()
	if err != nil {
		panic("failed to close database connection")
	}

	db = nil
}
