package handler

import (
	"github.com/lcardelli/fornecedores/config"
)

var (
	logger *config.Logger
)

// InitHandler initializes the handler
func InitHandler() {
	logger = config.GetLogger("handler")
	db = config.GetMysql()
}
