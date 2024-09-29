package handler

import (
	"github.com/lcardelli/fornecedores/config"
	"gorm.io/gorm"
)

var (
	logger *config.Logger
	db     *gorm.DB
)
// InitHandler initializes the handler
func InitHandler() {
	logger = config.GetLogger("handler")
	db = config.GetMysql()
}
