package config

import (
	"gorm.io/gorm"
	"errors"
)

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
	return errors.New("fake err")

	//return nil
}

func GetLogger(p string) *Logger {
	// Initialize Logger
	logger = NewLogger(p)
	return logger
}
