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
			auth.Use(handler.PermissionMiddleware("suppliers")).GET("/catalogo", handler.CatalogFornecedoresHandler)
			auth.GET("/lista-fornecedores", handler.ListaFornecedoresExternosHandler)

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
			auth.Use(handler.PermissionMiddleware("licenses")).GET("/licenses/view", handler.RenderViewLicensesPage)
			auth.GET("/licenses/list", handler.ListLicensesHandler)

			// Rotas de administração de fornecedores
			supplierAdmin := auth.Group("/")
			supplierAdmin.Use(handler.SupplierAdminMiddleware())
			{
				// Gerenciamento de Categorias e Serviços
				supplierAdmin.GET("/cadastro-categoria", handler.RenderCategoriaHandler)
				supplierAdmin.POST("/categories", handler.CreateCategoryHandler)
				supplierAdmin.PUT("/categories/:id", handler.UpdateCategoryHandler)
				supplierAdmin.DELETE("/categories/:id", handler.DeleteCategoryHandler)
				supplierAdmin.DELETE("/categories/batch", handler.DeleteMultipleCategories)
				
				supplierAdmin.GET("/services", handler.RenderServicePageHandler)
				supplierAdmin.POST("/services", handler.CreateServiceHandler)
				supplierAdmin.PUT("/services/:id", handler.UpdateServiceHandler)
				supplierAdmin.DELETE("/services/:id", handler.DeleteServiceHandler)
				supplierAdmin.DELETE("/services/batch", handler.DeleteMultipleServices)
				
				supplierAdmin.GET("/produtos", handler.RenderProductPageHandler)
				supplierAdmin.POST("/products", handler.CreateProductHandler)
				supplierAdmin.PUT("/products/:id", handler.UpdateProductHandler)
				supplierAdmin.DELETE("/products/:id", handler.DeleteProductHandler)
				supplierAdmin.DELETE("/products/batch", handler.DeleteMultipleProducts)
			}

			// Rotas de administração de licenças
			licenseAdmin := auth.Group("/")
			licenseAdmin.Use(handler.LicenseAdminMiddleware())
			{
				// Gerenciamento de Licenças
				licenseAdmin.GET("/licenses/manage", handler.RenderManageLicensesHandler)
				licenseAdmin.POST("/licenses", handler.CreateLicenseHandler)
				licenseAdmin.DELETE("/licenses/:id", handler.DeleteLicenseHandler)
				licenseAdmin.PUT("/licenses/:id", handler.UpdateLicenseHandler)
				
				// Gerenciamento de Softwares
				licenseAdmin.GET("/licenses/software", handler.RenderManageSoftwareHandler)
				licenseAdmin.POST("/licenses/software", handler.CreateSoftwareHandler)
				licenseAdmin.PUT("/licenses/software/:id", handler.UpdateSoftwareHandler)
				licenseAdmin.DELETE("/licenses/software/:id", handler.DeleteSoftwareHandler)
				licenseAdmin.GET("/licenses/software/:id", handler.GetSoftwareHandler)
				licenseAdmin.GET("/licenses/:id", handler.GetLicense)
			}

			// Mantenha o grupo admin apenas para gerenciamento de usuários
			adminUsers := auth.Group("/")
			adminUsers.Use(handler.AdminMiddleware())
			{
				// Gerenciamento de Usuários
				adminUsers.GET("/manage-users", handler.RenderManageUsersHandler)
				adminUsers.PUT("/users/:id/toggle-admin", handler.ToggleAdminHandler)
				adminUsers.DELETE("/users/:id", handler.DeleteUserHandler)
				adminUsers.GET("/users/:id/permissions", handler.GetUserPermissionsHandler)
				adminUsers.POST("/users/permissions", handler.UpdateUserPermissionsHandler)
			}
		}
	}
}
