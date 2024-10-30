package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func CreateProductHandler(c *gin.Context) {
	// Recebe o input do usuário
	var input struct {
		Name      string `json:"name" binding:"required"`
		ServiceID uint   `json:"service_id" binding:"required"`
	}

	// Valida o input do usuário
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove espaços do início e fim
	input.Name = strings.TrimSpace(input.Name)

	// Verifica se está vazio após o trim
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name cannot be empty"})
		return
	}

	// Verifica se o produto já existe
	if err := db.Where("name = ?", input.Name).First(&schemas.Product{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists"})
		return
	}

	// Verifica se o serviço existe
	var service schemas.Service
	if err := db.Where("id = ?", input.ServiceID).First(&service).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}

	// Cria o produto
	product := schemas.Product{
		Name:      input.Name,
		ServiceID: input.ServiceID,
	}

	// Salva o produto no banco de dados
	if err := db.Create(&product).Error; err != nil {
		log.Printf("Erro ao criar produto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Busca o produto criado com suas relações
	var createdProduct schemas.Product
	if err := db.Preload("Service").First(&createdProduct, product.ID).Error; err != nil {
		log.Printf("Erro ao buscar produto criado: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cria a resposta formatada
	response := schemas.ProductResponse{
		ID:        createdProduct.ID,
		Name:      createdProduct.Name,
		ServiceID: createdProduct.ServiceID,
		Service: schemas.ServiceResponse{
			ID:   createdProduct.Service.ID,
			Name: createdProduct.Service.Name,
		},
	}

	log.Printf("Produto criado com sucesso: ID %d", product.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Produto criado com sucesso",
		"product": response,
	})
}
