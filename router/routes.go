package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/handler"
)

// InitializeRoutes initializes the routes for the application
func InitializeRoutes(router *gin.Engine) {
	handler.InitHandler()
	// Create a new group for the v1 API
	v1 := router.Group("/api/v1")
	{

		v1.GET("/supplier", handler.ShowOpeningHandler)

		v1.POST("/suppliers", handler.CreateSupplierHandler)

		v1.DELETE("/suppliers", handler.DeleteOpeningHandler)

		v1.PUT("/suppliers", handler.UpdateOpeningHandler)

		v1.GET("/suppliers", handler.ListOpeningHandler)
	}
}
