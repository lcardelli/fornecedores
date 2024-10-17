package handler

import (
	"log"
	"net/http"
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
	log.Printf("Buscando serviços para a categoria ID: %s", categoryID)

	var services []schemas.Service
	if err := db.Where("category_id = ?", categoryID).Find(&services).Error; err != nil {
		log.Printf("Erro ao buscar serviços: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	log.Printf("Encontrados %d serviços para a categoria %s", len(services), categoryID)
	c.JSON(http.StatusOK, services)
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
