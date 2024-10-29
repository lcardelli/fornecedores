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
			auth.GET("/dashboard", handler.DashboardHandler)                          // Renderiza a página do dashboard
			auth.GET("/catalogo", handler.CatalogFornecedoresHandler)                 // Renderiza a página do catálogo de fornecedores
			auth.GET("/lista-fornecedores", handler.ListaFornecedoresExternosHandler) // Renderiza a página da lista de fornecedores externos
			auth.GET("/cadastro-fornecedor", handler.FormRegisterHandler)             // Renderiza a página de cadastro de fornecedor
			auth.GET("/services", handler.RenderServicePageHandler)                   // Renderiza a página de cadastro de serviços
			auth.GET("/cadastro-categoria", handler.RenderCategoriaHandler)           // Renderiza a página de cadastro de categoria

			// Rotas para fornecedores
			auth.POST("/suppliers", handler.CreateSupplierHandler)       // Cria um novo fornecedor
			auth.GET("/suppliers", handler.ListSupplierHandler)          // Lista todos os fornecedores
			auth.GET("/suppliers/:id", handler.ShowSupplierHandler)      // Mostra um fornecedor pelo ID
			auth.PUT("/suppliers/:id", handler.UpdateSupplierHandler)    // Atualiza um fornecedor pelo ID
			auth.DELETE("/suppliers/:id", handler.DeleteSupplierHandler) // Deleta um fornecedor pelo ID

			// Logout do Google
			auth.GET("/auth/google/logout", handler.GoogleLogout)

			// Rotas para categorias
			auth.POST("/categories", handler.CreateCategoryHandler)       // Cria uma nova categoria
			auth.GET("/categories", handler.ListCategoriesHandler)        // Lista todas as categorias
			auth.PUT("/categories/:id", handler.UpdateCategoryHandler)    // Atualiza uma categoria pelo ID
			auth.DELETE("/categories/:id", handler.DeleteCategoryHandler) // Deleta uma categoria pelo ID

			auth.POST("/services", handler.CreateServiceHandler)                                // Cria um novo serviço
			auth.GET("/service-list", handler.ListServicesHandler)                              // Lista todos os serviços
			auth.PUT("/services/:id", handler.UpdateServiceHandler)                             // Atualiza um serviço pelo ID
			auth.DELETE("/services/:id", handler.DeleteServiceHandler)                          // Deleta um serviço pelo ID
			auth.GET("/services-by-category/:id", handler.GetServicesByCategoryHandler) // busca os serviços por categoria
			auth.GET("/suppliers-by-id", handler.GetSupplierHandler)                            // busca os fornecedores pelo ID

			// Adicione esta nova rota ao seu router dentro do grupo auth
			auth.DELETE("/categories/batch", handler.DeleteMultipleCategories)

			// Adicione esta nova rota ao seu router
			auth.DELETE("/services/batch", handler.DeleteMultipleServices)

			// Rotas para produtos
			auth.GET("/products", handler.GetProductsHandler)
			auth.DELETE("/products/batch", handler.DeleteMultipleProducts)
			auth.POST("/products", handler.CreateProductHandler)
			auth.PUT("/products/:id", handler.UpdateProductHandler)
			auth.DELETE("/products/:id", handler.DeleteProductHandler)
			auth.GET("/products-by-service/:id", handler.GetProductsByServiceHandler)
		}
	}
}
