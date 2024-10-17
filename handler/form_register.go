package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func FormRegisterHandler(c *gin.Context) {
	user, _ := c.Get("user")
	typedUser := user.(schemas.User)

	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Não vamos mais buscar todos os serviços aqui
	// Vamos buscar os serviços dinamicamente via AJAX

	fornecedores, err := getFornecedoresFromDatabase()
	if err != nil {
		log.Printf("Erro ao buscar fornecedores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados de fornecedores: " + err.Error()})
		return
	}

	// Filtrar fornecedores com base nos parâmetros de busca
	search := c.Query("search")
	name := c.Query("name")
	cnpj := c.Query("cnpj")

	filteredFornecedores := filterFornecedores(fornecedores, search, name, cnpj)

	log.Printf("Número de fornecedores filtrados: %d", len(filteredFornecedores))

	c.HTML(http.StatusOK, "form_register.html", gin.H{
		"user":         typedUser,
		"Categories":   categories,
		"Fornecedores": filteredFornecedores,
		"activeMenu":   "cadastro-fornecedor",
		"search":       search,
		"name":         name,
		"cnpj":         cnpj,
	})
}

// Adicione esta nova função para buscar serviços por categoria
func GetServicesByCategoryHandler(c *gin.Context) {
	categoryID := c.Param("categoryId")
	categoryIDInt, _ := strconv.Atoi(categoryID)

	var services []schemas.Service
	if err := db.Joins("JOIN supplier_services ON services.id = supplier_services.service_id").
		Joins("JOIN supplier_links ON supplier_links.id = supplier_services.supplier_link_id").
		Where("supplier_links.category_id = ?", categoryIDInt).
		Distinct().Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	log.Printf("Serviços encontrados para a categoria %d: %+v", categoryIDInt, services)

	// Convertendo para ServiceResponse para garantir que os campos JSON estejam corretos
	var serviceResponses []schemas.ServiceResponse
	for _, service := range services {
		serviceResponses = append(serviceResponses, schemas.ServiceResponse{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
			Price:       service.Price,
		})
	}

	c.JSON(http.StatusOK, serviceResponses)
}

func filterFornecedores(fornecedores []Fornecedor, search, name, cnpj string) []Fornecedor {
	var filtered []Fornecedor

	for _, f := range fornecedores {
		if (search == "" || (strings.Contains(strings.ToLower(f.NOME.String), strings.ToLower(search)) ||
			strings.Contains(f.CGCCFO.String, search))) &&
			(name == "" || strings.Contains(strings.ToLower(f.NOME.String), strings.ToLower(name))) &&
			(cnpj == "" || strings.Contains(f.CGCCFO.String, cnpj)) {
			filtered = append(filtered, f)
		}
	}

	return filtered
}
