package router

import "github.com/gin-gonic/gin"

func Initialize() {
	// Create a new Gin router
	router := gin.Default()
	// Initialize the routes
	InitializeRoutes(router)
	// Start the Gin server on port 8080
	router.Run(":8080")
}
