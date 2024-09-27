package main

import (

	"github.com/lcardelli/fornecedores/config"
	"github.com/lcardelli/fornecedores/router"
)

var(
	logger *config.Logger
)


func main() {

	logger = config.GetLogger("main")

	// Initialize configs
	err := config.Init()
	if err != nil {
		logger.Errorf("config initialization error: %v", err)
		return
	}

	// Initialize the router
	router.Initialize()
}
