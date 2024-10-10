package handler

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func CatalogFornecedoresHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID") // Obtém o ID do usuário da sessão

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user schemas.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Obter parâmetros de filtro
	categoryID := c.Query("category")
	serviceID := c.Query("service")
	supplierName := c.Query("name")

	// Buscar categorias para o filtro
	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categorias"})
		return
	}

	// Buscar serviços para o filtro
	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	// Construir a query para fornecedores
	query := db.Preload("Category").Preload("Services").Preload("Services.Service")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if serviceID != "" {
		query = query.Joins("JOIN supplier_services ON suppliers.id = supplier_services.supplier_id").
			Where("supplier_services.service_id = ?", serviceID)
	}
	if supplierName != "" {
		query = query.Where("name LIKE ?", "%"+supplierName+"%")
	}

	var suppliers []schemas.Supplier
	if err := query.Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar fornecedores"})
		return
	}

	// Carregar templates
	tmpl := template.Must(template.ParseGlob("templates/*"))
	
	// Renderizar o template catalogo.html
	err := tmpl.ExecuteTemplate(c.Writer, "catalogo.html", gin.H{
		"user":       user,
		"suppliers":  suppliers,
		"categories": categories,
		"services":   services,
		"filters": gin.H{
			"category": categoryID,
			"service":  serviceID,
			"name":     supplierName,
		},
		"activeMenu": "catalogo",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao renderizar o template"})
		return
	}
}
