package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// CatalogFornecedoresHandler exibe o catálogo de fornecedores
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
			Where("supplier_links.deleted_at IS NULL"). // Adiciona esta linha para excluir registros deletados
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
	query := db.Preload("Category").Preload("Services", "supplier_services.deleted_at IS NULL").Preload("Services.Service")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if serviceID != "" {
		query = query.Joins("JOIN supplier_services ON supplier_links.id = supplier_services.supplier_link_id").
			Where("supplier_services.service_id = ? AND supplier_services.deleted_at IS NULL", serviceID)
	}

	var supplierLinks []schemas.SupplierLink
	if err := query.Where("supplier_links.deleted_at IS NULL").Find(&supplierLinks).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar vínculos de fornecedores"})
		return
	}

	// Buscar informações externas dos fornecedores
	fornecedores, err := getFornecedoresExternosFromDatabase()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar dados de fornecedores externos"})
		return
	}

	// Modificar a estrutura para armazenar fornecedores com informações adicionais
	type FornecedorComDetalhes struct {
		Fornecedor
		ID        uint // ID do SupplierLink
		Categoria string
		Servicos  []string
		CNPJ      string // Adicionando CNPJ para facilitar a identificação
	}

	// Modificar a lógica de filtragem e criação da lista de fornecedores
	var fornecedoresComDetalhes []FornecedorComDetalhes
	for _, f := range fornecedores {
		for _, link := range supplierLinks {
			if f.CGCCFO.String == link.CNPJ {
				detalhe := FornecedorComDetalhes{
					Fornecedor: f,
					ID:         link.ID,
					CNPJ:       link.CNPJ,
					Categoria:  link.Category.Name,
					Servicos:   make([]string, 0),
				}
				for _, s := range link.Services {
					if s.DeletedAt.Time.IsZero() { // Verifica se o serviço não foi deletado
						detalhe.Servicos = append(detalhe.Servicos, s.Service.Name)
					}
				}
				// Adiciona o fornecedor apenas se ele tiver serviços ativos
				if len(detalhe.Servicos) > 0 {
					fornecedoresComDetalhes = append(fornecedoresComDetalhes, detalhe)
				}
				break
			}
		}
	}

	// Aplicar filtro de nome do fornecedor, se fornecido
	if supplierName != "" {
		var filteredFornecedores []FornecedorComDetalhes
		for _, f := range fornecedoresComDetalhes {
			if strings.Contains(strings.ToLower(f.NOME.String), strings.ToLower(supplierName)) {
				filteredFornecedores = append(filteredFornecedores, f)
			}
		}
		fornecedoresComDetalhes = filteredFornecedores
	}

	// Renderizar o template catalogo.html
	c.HTML(http.StatusOK, "catalogo.html", gin.H{
		"user":       user,
		"suppliers":  fornecedoresComDetalhes,
		"categories": categories,
		"services":   services,
		"filters": gin.H{
			"category":   categoryID,
			"service":    serviceID,
			"name":       supplierName,
			"activeMenu": "catalogo",
		},
	})
}

// GetCategoriesHandler busca todas as categorias
func GetCategoriesHandler(c *gin.Context) {
	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categorias"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetServicesHandler busca todos os serviços
func GetServicesHandler(c *gin.Context) {
	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}
	c.JSON(http.StatusOK, services)
}

// GetSupplierHandler busca um fornecedor pelo ID
func GetSupplierHandler(c *gin.Context) {
	supplierID := c.Query("id")
	if supplierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do fornecedor não fornecido"})
		return
	}

	var supplierLink schemas.SupplierLink
	if err := db.Preload("Category").
		Preload("Services", "deleted_at IS NULL").
		Preload("Services.Service").
		First(&supplierLink, supplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fornecedor não encontrado"})
		return
	}

	// Criar uma estrutura que combina as informações do SupplierLink
	response := gin.H{
		"ID":       supplierLink.ID,
		"CNPJ":     supplierLink.CNPJ,
		"Category": supplierLink.Category,
		"Services": supplierLink.Services,
		// ... outras informações relevantes ...
	}

	c.JSON(http.StatusOK, response)
}
