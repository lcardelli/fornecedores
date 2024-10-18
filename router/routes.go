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
			auth.GET("/services-by-category/:categoryId", handler.GetServicesByCategoryHandler)
			auth.POST("/suppliers", handler.CreateSupplierHandler)
			auth.GET("/suppliers", handler.ListSupplierHandler)
			auth.GET("/suppliers/:id", handler.ShowSupplierHandler)
			auth.PUT("/suppliers/:id", handler.UpdateSupplierHandler)
			auth.DELETE("/suppliers/:id", handler.DeleteSupplierHandler)

			auth.GET("/auth/google/logout", handler.GoogleLogout)
			auth.GET("/cadastro-categoria", handler.CadastroCategoriaHandler)
			auth.POST("/categories", handler.CreateCategoryHandler)
			auth.GET("/categories", handler.ListCategoriesHandler)
			auth.PUT("/categories/:id", handler.UpdateCategoryHandler)
			auth.DELETE("/categories/:id", handler.DeleteCategoryHandler)
			auth.GET("/services", handler.CadastroServicoHandler)
			auth.POST("/services", handler.CreateServiceHandler)
			auth.GET("/service-list", handler.ListServicesHandler) // Nova rota para listar serviços
			auth.PUT("/services/:id", handler.UpdateServiceHandler)
			auth.DELETE("/services/:id", handler.DeleteServiceHandler)

		}
	}
}
