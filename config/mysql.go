package config

import (
	"os"

	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

// InitializeMysql initializes the MySQL database and performs auto-migration
func InitializeMysql() (*gorm.DB, error) {
	logger := GetLogger("mysql")

	// Load the environment variables
	_ = godotenv.Load() 
	dsn := os.Getenv("DATABASE_URL")
	// Create the connection with the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	// AutoMigrate
	err = db.AutoMigrate(&schemas.Supplier{}, &schemas.SupplierCategory{}, &schemas.SupplierService{}, &schemas.User{})
	if err != nil {
		logger.Errorf("Failed to migrate database: %v", err)
		return nil, err
	}

	// Return the database connection
	return db, nil
}