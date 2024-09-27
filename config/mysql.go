package config

import (
	"os"

	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMysql() (*gorm.DB, error) {
	logger := GetLogger("mysql")
	dbpath := "./db/main.db"
	//Check if the database is already connected
	_, err := os.Stat(dbpath)
	if os.IsNotExist(err) {
		logger.Info("Database not found, creating...")
		err = os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			logger.Errorf("Failed to create database directory: %v", err)
			return nil, err
		}
		file, err := os.Create(dbpath)
		if err != nil {
			logger.Errorf("Failed to create database file: %v", err)
			return nil, err
		}
		file.Close()
	}

	// Initialize MySQL
	db, err := gorm.Open(mysql.Open(dbpath), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	// AutoMigrate
	err = db.AutoMigrate(&schemas.Supplier{})
	if err != nil {
		logger.Errorf("Failed to migrate database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&schemas.SupplierCategory{})
	if err != nil {
		logger.Errorf("Failed to migrate SupplierCategory: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&schemas.SupplierService{})
	if err != nil {
		logger.Errorf("Failed to migrate SupplierService: %v", err)
		return nil, err
	}
	// Return the database connection
	return db, nil
} 
