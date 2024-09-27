package config

import (
	"fmt"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
	var err error
	db, err = InitializeMysql()
	if err != nil {
		return fmt.Errorf("failed to initialize mysql: %w", err)
	}

	return nil
}

func GetMysql() *gorm.DB {
	return db
}

func GetLogger(p string) *Logger {
	// Initialize Logger
	logger = NewLogger(p)
	return logger
}
