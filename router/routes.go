package router

import (
	"github.com/gin-gonic/gin"
	docs "github.com/lcardelli/fornecedores/docs"
	"github.com/lcardelli/fornecedores/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

		v1.GET("/auth/google", handler.GoogleLogin)
		v1.GET("/auth/google/callback", handler.GoogleCallback)
		v1.GET("/index", handler.IndexHandler)

		v1.GET("/dashboard", handler.DashboardHandler) // Adicionando a rota do dashboard

		// Rotas protegidas
		v1.Use(handler.AuthMiddleware()) // Aplicando o middleware
	}

	// Initializei Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
