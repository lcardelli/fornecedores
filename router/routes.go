package router

import "github.com/gin-gonic/gin"

// InitializeRoutes initializes the routes for the application
func InitializeRoutes(router *gin.Engine) {
	// Create a new group for the v1 API
	v1 := router.Group("/api/v1")
	{
		// Create a new group for the suppliers
		v1.GET("/supplie", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg": "Suppliers",
			})
		})

		v1.POST("/suppliers", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg": "Suppliers",
			})
		})

		v1.DELETE("/suppliers", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg": "Suppliers",
			})
		})

		v1.PUT("/suppliers", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg": "Suppliers",
			})
		})

		v1.PUT("/suppliers", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg": "Suppliers",
			})
		})
	}
}