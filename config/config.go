package config

import (
	"fmt"
	"gorm.io/gorm"
)

// Initialize the database and logger
var (
	db     *gorm.DB
	logger *Logger
)
// Initialize the database and logger
func Init() error {
	var err error
	db, err = InitializeMysql()
	if err != nil {
		return fmt.Errorf("failed to initialize mysql: %w", err)
	}

	return nil
}

// Get the database
func GetMysql() *gorm.DB {
	return db
}

// Get the logger
func GetLogger(p string) *Logger {
	// Initialize Logger
	logger = NewLogger(p)
	return logger
}
