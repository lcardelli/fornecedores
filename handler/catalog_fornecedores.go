package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func CatalogFornecedoresHandler(c *gin.Context) {
	// Obter o usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Usuário não autenticado"})
		return
	}
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao obter informações do usuário"})
		return
	}

	// Obter parâmetros de filtro
	categoryID := c.Query("category")
	serviceID := c.Query("service")
	supplierName := c.Query("name")

	// Buscar categorias para o filtro
	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar categorias"})
		return
	}

	// Buscar serviços para o filtro
	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	// Construir a query para SupplierLinks
	query := db.Preload("Category").Preload("Services").Preload("Services.Service")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if serviceID != "" {
		query = query.Joins("JOIN supplier_services ON supplier_links.id = supplier_services.supplier_link_id").
			Where("supplier_services.service_id = ?", serviceID)
	}

	var supplierLinks []schemas.SupplierLink
	if err := query.Find(&supplierLinks).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar vínculos de fornecedores"})
		return
	}

	// Buscar informações externas dos fornecedores
	fornecedores, err := getFornecedoresFromDatabase()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar dados de fornecedores externos"})
		return
	}

	// Filtrar fornecedores com base nos critérios
	var filteredFornecedores []Fornecedor
	for _, f := range fornecedores {
		if supplierName != "" && (f.NOME.String == "" || f.NOME.String != supplierName) {
			continue
		}
		for _, link := range supplierLinks {
			if f.CGCCFO.String == link.CNPJ {
				filteredFornecedores = append(filteredFornecedores, f)
				break
			}
		}
	}

	// Renderizar o template catalogo.html
	c.HTML(http.StatusOK, "catalogo.html", gin.H{
		"user":         user,
		"suppliers":    filteredFornecedores, // Mudamos de "fornecedores" para "suppliers"
		"categories":   categories,
		"services":     services,
		"filters": gin.H{
			"category": categoryID,
			"service":  serviceID,
			"name":     supplierName,
			"activeMenu": "catalogo",
		},
	})
}
