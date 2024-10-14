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
		// Rotas públicas
		v1.GET("/auth/google", handler.GoogleLogin)
		v1.GET("/auth/google/callback", handler.GoogleCallback)
		v1.GET("/index", handler.IndexHandler) // Rota para a página de login

		// Swagger
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Rotas protegidas
		auth := v1.Group("/")
		auth.Use(handler.AuthMiddleware())
		{
			auth.GET("/dashboard", handler.DashboardHandler)
			auth.GET("/catalogo", handler.CatalogFornecedoresHandler)
			auth.GET("/lista-fornecedores", handler.ListaFornecedoresHandler)
			auth.GET("/cadastro-fornecedor", handler.FormRegisterHandler)
			auth.GET("/services", handler.GetServicesByCategoryHandler)

			auth.GET("/supplier", handler.ShowSupplierHandler)
			auth.POST("/suppliers", handler.CreateSupplierHandler)
			auth.DELETE("/suppliers", handler.DeleteSupplierHandler)
			auth.PUT("/suppliers", handler.UpdateSupplierHandler)
			auth.GET("/suppliers", handler.ListSupplierHandler)
			auth.GET("/suppliers/:id/services", handler.ListServicesHandler)

			auth.GET("/auth/google/logout", handler.GoogleLogout)
		}
	}
}
