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
	if categoryID != "" {
		categoryIDInt, _ := strconv.Atoi(categoryID)
		if err := db.Joins("JOIN supplier_services ON services.id = supplier_services.service_id").
			Joins("JOIN supplier_links ON supplier_links.id = supplier_services.supplier_link_id").
			Where("supplier_links.category_id = ?", categoryIDInt).
			Distinct().Find(&services).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar serviços"})
			return
		}
	} else {
		if err := db.Find(&services).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar serviços"})
			return
		}
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

	// Modificar a estrutura para armazenar fornecedores com informações adicionais
	type FornecedorComDetalhes struct {
		Fornecedor
		Categoria string
		Servicos  []string
	}

	// Modificar a lógica de filtragem e criação da lista de fornecedores
	var fornecedoresComDetalhes []FornecedorComDetalhes
	for _, f := range fornecedores {
		for _, link := range supplierLinks {
			if f.CGCCFO.String == link.CNPJ {
				detalhe := FornecedorComDetalhes{
					Fornecedor: f,
					Categoria:  link.Category.Name,
					Servicos:    make([]string, 0),
				}
				for _, s := range link.Services {
					detalhe.Servicos = append(detalhe.Servicos, s.Service.Name)
				}
				fornecedoresComDetalhes = append(fornecedoresComDetalhes, detalhe)
				break
			}
		}
	}

	// Renderizar o template catalogo.html
	c.HTML(http.StatusOK, "catalogo.html", gin.H{
		"user":         user,
		"suppliers":    fornecedoresComDetalhes, // Agora usando fornecedoresComDetalhes
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



