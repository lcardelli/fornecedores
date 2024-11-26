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
			// Rotas públicas (qualquer usuário autenticado)
			auth.GET("/dashboard", handler.DashboardHandler)
			auth.GET("/catalogo", handler.CatalogFornecedoresHandler)
			auth.GET("/lista-fornecedores", handler.ListaFornecedoresExternosHandler)

			// Rotas que requerem privilégios de administrador
			admin := auth.Group("/")
			admin.Use(handler.AdminMiddleware())
			{
				// Gerenciamento de Áreas
				admin.GET("/cadastro-categoria", handler.RenderCategoriaHandler)
				admin.POST("/categories", handler.CreateCategoryHandler)
				admin.PUT("/categories/:id", handler.UpdateCategoryHandler)
				admin.DELETE("/categories/:id", handler.DeleteCategoryHandler)
				admin.DELETE("/categories/batch", handler.DeleteMultipleCategories)

				// Gerenciamento de Categorias
				admin.GET("/services", handler.RenderServicePageHandler)
				admin.POST("/services", handler.CreateServiceHandler)
				admin.PUT("/services/:id", handler.UpdateServiceHandler)
				admin.DELETE("/services/:id", handler.DeleteServiceHandler)
				admin.DELETE("/services/batch", handler.DeleteMultipleServices)

				// Gerenciamento de Produtos
				admin.GET("/produtos", handler.RenderProductPageHandler)
				admin.POST("/products", handler.CreateProductHandler)
				admin.PUT("/products/:id", handler.UpdateProductHandler)
				admin.DELETE("/products/:id", handler.DeleteProductHandler)
				admin.DELETE("/products/batch", handler.DeleteMultipleProducts)

				// Gerenciamento de Usuários
				admin.GET("/manage-users", handler.RenderManageUsersHandler)
				admin.PUT("/users/:id/toggle-admin", handler.ToggleAdminHandler)
				admin.DELETE("/users/:id", handler.DeleteUserHandler)

				// Gerenciamento de Licenças
				admin.GET("/licenses/manage", handler.RenderManageLicensesHandler)
				admin.POST("/licenses", handler.CreateLicenseHandler)
				admin.DELETE("/licenses/:id", handler.DeleteLicenseHandler)
				admin.PUT("/licenses/:id", handler.UpdateLicenseHandler)

				// Gerenciamento de Softwares
				admin.GET("/licenses/software", handler.RenderManageSoftwareHandler)
				admin.POST("/licenses/software", handler.CreateSoftwareHandler)
				admin.PUT("/licenses/software/:id", handler.UpdateSoftwareHandler)
				admin.DELETE("/licenses/software/:id", handler.DeleteSoftwareHandler)
				admin.GET("/licenses/software/:id", handler.GetSoftwareHandler)
				admin.GET("/licenses/:id", handler.GetLicense)

				// Dentro do grupo admin
				admin.GET("/users/:id/permissions", handler.GetUserPermissionsHandler)
				admin.POST("/users/permissions", handler.UpdateUserPermissionsHandler)
			}

			// Rotas para fornecedores
			auth.POST("/suppliers", handler.CreateSupplierHandler)       // Cria um novo fornecedor
			auth.GET("/suppliers", handler.ListSupplierHandler)          // Lista todos os fornecedores
			auth.GET("/suppliers/:id", handler.ShowSupplierHandler)      // Mostra um fornecedor pelo ID
			auth.PUT("/suppliers/:id", handler.UpdateSupplierHandler)    // Atualiza um fornecedor pelo ID
			auth.DELETE("/suppliers/:id", handler.DeleteSupplierHandler) // Deleta um fornecedor pelo ID

			// Logout do Google
			auth.GET("/auth/google/logout", handler.GoogleLogout)

			// Rotas para categorias
			auth.GET("/categories", handler.ListCategoriesHandler)                      // Lista todas as categorias
			auth.GET("/services-by-category/:id", handler.GetServicesByCategoryHandler) // busca os serviços por categoria
			auth.GET("/suppliers-by-id", handler.GetSupplierHandler)                    // busca os fornecedores pelo ID

			// Rotas para produtos
			auth.GET("/products-list", handler.ListSupplierProducts)                  // Lista todos os produtos
			auth.GET("/products-by-service/:id", handler.GetProductsByServiceHandler) // Busca produtos por serviço
			auth.GET("/products", handler.GetProductsHandler)                         // Busca todos os produtos

			// Adicione esta rota para listar serviços
			auth.GET("/service-list", handler.ListServicesHandler) // Nova rota para listar serviços

			// Adicione estas rotas junto com as outras rotas existentes
			auth.GET("/licenses/view", handler.RenderViewLicensesPage)
			auth.GET("/licenses/list", handler.ListLicensesHandler)

		}
	}
}
