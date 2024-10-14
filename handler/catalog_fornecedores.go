package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func CatalogFornecedoresHandler(c *gin.Context) {
	// Obter o usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do usuário"})
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
	if categoryID != "" {
		categoryIDInt, _ := strconv.Atoi(categoryID)
		if err := db.Joins("JOIN supplier_services ON services.id = supplier_services.service_id").
			Joins("JOIN suppliers ON suppliers.id = supplier_services.supplier_id").
			Where("suppliers.category_id = ?", categoryIDInt).
			Distinct().Find(&services).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
			return
		}
	} else {
		if err := db.Find(&services).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
			return
		}
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

	// Renderizar o template catalogo.html
	c.HTML(http.StatusOK, "catalogo.html", gin.H{
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
}
