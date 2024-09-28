package config

import (
	"os"

	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

func InitializeMysql() (*gorm.DB, error) {
	logger := GetLogger("mysql")

	// Inicializa o MySQL
	_ = godotenv.Load() // Carrega as variáveis do arquivo .env
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	// AutoMigrate
	err = db.AutoMigrate(&schemas.Supplier{}, &schemas.SupplierCategory{}, &schemas.SupplierService{})
	if err != nil {
		logger.Errorf("Failed to migrate database: %v", err)
		return nil, err
	}

	// Retorna a conexão com o banco de dados
	return db, nil
}
