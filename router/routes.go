package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/handler"
	docs "github.com/lcardelli/fornecedores/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// InitializeRoutes initializes the routes for the application
func InitializeRoutes(router *gin.Engine) {
	handler.InitHandler()
	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	// Create a new group for the v1 API
	v1 := router.Group(basePath)
	{

		v1.GET("/supplier", handler.ShowSupplierHandler)

		v1.POST("/suppliers", handler.CreateSupplierHandler)

		v1.DELETE("/suppliers", handler.DeleteSupplierHandler)

		v1.PUT("/suppliers", handler.UpdateSupplierHandler)

		v1.GET("/suppliers", handler.ListSupplierHandler)
	}

	// Initialize Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
