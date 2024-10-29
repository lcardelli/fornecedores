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
	productID := c.Query("product")
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

	// Buscar produtos para o filtro
	var products []schemas.Product
	if serviceID != "" {
		serviceIDInt, _ := strconv.Atoi(serviceID)
		if err := db.Where("service_id = ?", serviceIDInt).Find(&products).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar produtos"})
			return
		}
	} else {
		if err := db.Find(&products).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar produtos"})
			return
		}
	}

	// Construir a query para SupplierLinks
	query := db.Preload("Category").
		Preload("Services", "supplier_services.deleted_at IS NULL").
		Preload("Services.Service").
		Preload("Products", "supplier_products.deleted_at IS NULL").
		Preload("Products.Product")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if serviceID != "" {
		query = query.Joins("JOIN supplier_services ON supplier_links.id = supplier_services.supplier_link_id").
			Where("supplier_services.service_id = ? AND supplier_services.deleted_at IS NULL", serviceID)
	}
	if productID != "" {
		productIDInt, _ := strconv.Atoi(productID)
		query = query.Joins("JOIN supplier_products ON supplier_links.id = supplier_products.supplier_link_id").
			Where("supplier_products.product_id = ? AND supplier_products.deleted_at IS NULL", productIDInt)
	}

	// Buscar vínculos de fornecedores
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
		Produtos  []string
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
					Produtos:   make([]string, 0),
				}
				for _, s := range link.Services {
					if s.DeletedAt.Time.IsZero() { // Verifica se o serviço não foi deletado
						detalhe.Servicos = append(detalhe.Servicos, s.Service.Name)
					}
				}
				for _, p := range link.Products {
					if p.DeletedAt.Time.IsZero() { // Verifica se o produto não foi deletado
						detalhe.Produtos = append(detalhe.Produtos, p.Product.Name)
					}
				}
				// Se há filtro de produto, verifica se o fornecedor tem o produto específico
				if productID != "" {
					hasProduct := false
					for _, p := range link.Products {
						if p.DeletedAt.Time.IsZero() && strconv.Itoa(int(p.ProductID)) == productID {
							hasProduct = true
							break
						}
					}
					if !hasProduct {
						continue
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
		"products":   products,
		"filters": gin.H{
			"category":   categoryID,
			"service":    serviceID,
			"product":    productID,
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

// GetProductsHandler busca todos os produtos
func GetProductsHandler(c *gin.Context) {
	var products []schemas.Product
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produtos"})
		return
	}
	c.JSON(http.StatusOK, products)
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
