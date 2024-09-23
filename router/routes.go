package router

import "github.com/gin-gonic/gin"

func Initialize() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route for the root path
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the Gin server on port 8080
	router.Run(":8080")
}
