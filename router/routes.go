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

	v1 := router.Group(basePath)
	{
		// Rotas públicas
		v1.GET("/auth/google", handler.GoogleLogin)
		v1.GET("/auth/google/callback", handler.GoogleCallback)
		v1.GET("/index", handler.IndexHandler)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Rotas protegidas
		auth := v1.Group("/")
		auth.Use(handler.AuthMiddleware())
		{
			// Rotas públicas (qualquer usuário autenticado)
			auth.GET("/dashboard", handler.DashboardHandler)
			auth.GET("/auth/google/logout", handler.GoogleLogout)

			// Rotas de fornecedores - visualização
			fornecedores := auth.Group("/")
			fornecedores.Use(handler.PermissionMiddleware("suppliers"))
			{
				fornecedores.GET("/catalogo", handler.CatalogFornecedoresHandler)
				fornecedores.GET("/lista-fornecedores", handler.ListaFornecedoresExternosHandler)
				fornecedores.GET("/suppliers", handler.ListSupplierHandler)
				fornecedores.GET("/suppliers/:id", handler.ShowSupplierHandler)
				fornecedores.GET("/categories", handler.ListCategoriesHandler)
				fornecedores.GET("/services-by-category/:id", handler.GetServicesByCategoryHandler)
				fornecedores.GET("/suppliers-by-id", handler.GetSupplierHandler)
				fornecedores.GET("/products-list", handler.ListSupplierProducts)
				fornecedores.GET("/products-by-service/:id", handler.GetProductsByServiceHandler)
				fornecedores.GET("/products", handler.GetProductsHandler)
				fornecedores.GET("/service-list", handler.ListServicesHandler)
			}

			// Rotas de fornecedores - administração
			supplierAdmin := auth.Group("/")
			supplierAdmin.Use(handler.PermissionMiddleware("supplier_admin"))
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

				supplierAdmin.POST("/suppliers", handler.CreateSupplierHandler)
				supplierAdmin.PUT("/suppliers/:id", handler.UpdateSupplierHandler)
				supplierAdmin.DELETE("/suppliers/:id", handler.DeleteSupplierHandler)
			}

			// Rotas de licenças - visualização
			licencas := auth.Group("/")
			licencas.Use(handler.PermissionMiddleware("licenses"))
			{
				licencas.GET("/licenses/view", handler.RenderViewLicensesPage)
				licencas.GET("/licenses/list", handler.ListLicensesHandler)
			}

			// Rotas de licenças - administração
			licenseAdmin := auth.Group("/")
			licenseAdmin.Use(handler.PermissionMiddleware("license_admin"))
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
				licenseAdmin.DELETE("/licenses/batch", handler.DeleteBatchLicenses)
			}

			// Rotas de administração global (apenas admin global)
			globalAdmin := auth.Group("/")
			globalAdmin.Use(handler.GlobalAdminMiddleware())
			{
				globalAdmin.GET("/manage-users", handler.RenderManageUsersHandler)
				globalAdmin.PUT("/users/:id/toggle-admin", handler.ToggleAdminHandler)
				globalAdmin.DELETE("/users/:id", handler.DeleteUserHandler)
				globalAdmin.GET("/users/:id/permissions", handler.GetUserPermissionsHandler)
				globalAdmin.POST("/users/permissions", handler.UpdateUserPermissionsHandler)
			}

			// Rotas de departamentos
			api := auth.Group("/")
			api.Use(handler.PermissionMiddleware("api"))
			{
				api.GET("/departments", handler.GetDepartmentsHandler)
			}

			// Rotas de Contratos
			contracts := auth.Group("/contracts")
			contracts.Use(handler.AuthMiddleware())

			// Rotas de visualização
			contracts.GET("/manage", handler.ContractAdminMiddleware(), handler.RenderManageContractsHandler) 
			contracts.GET("/list", handler.ContractViewMiddleware(), handler.RenderListContractsHandler)
			contracts.GET("", handler.ContractViewMiddleware(), handler.GetAllContractsHandler)
			contracts.GET("/:id", handler.ContractViewMiddleware(), handler.GetContractHandler)
			contracts.GET("/:id/aditivos", handler.GetContractAditivosHandler)

			// Rotas que requerem permissão de administração
			contractsAdmin := contracts.Group("")
			contractsAdmin.Use(handler.PermissionMiddleware("contract_admin"))
			{
				contractsAdmin.POST("", handler.CreateContractHandler)
				contractsAdmin.PUT("/:id", handler.UpdateContractHandler)
				contractsAdmin.DELETE("/:id", handler.DeleteContractHandler)
				contractsAdmin.POST("/batch-delete", handler.DeleteBatchContracts)
				contractsAdmin.GET("/download/:id", handler.DownloadContractAttachmentHandler)

				// Rotas de aditivos
				contractsAdmin.POST("/aditivos", handler.CreateContractAditivoHandler)
				contractsAdmin.PUT("/aditivos/:id", handler.UpdateContractAditivoHandler)
				contractsAdmin.DELETE("/aditivos/:id", handler.DeleteContractAditivoHandler)
			}
		}
	}
}
