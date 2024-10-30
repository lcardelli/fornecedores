package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// CreateServiceHandler lida com a criação de um novo serviço
func CreateServiceHandler(c *gin.Context) {
	log.Println("Iniciando CreateServiceHandler")

	var input struct {
		Name       string `json:"name" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove espaços do início e fim
	input.Name = strings.TrimSpace(input.Name)

	// Verifica se está vazio após o trim
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service name cannot be empty"})
		return
	}

	// Verifica se o serviço já existe
	if err := db.Where("name = ?", input.Name).First(&schemas.Service{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service already exists"})
		return
	}

	// Verifica se a categoria existe
	var category schemas.SupplierCategory
	if err := db.First(&category, input.CategoryID).Error; err != nil {
		log.Printf("Categoria não encontrada: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	service := schemas.Service{
		Name:       input.Name,
		CategoryID: input.CategoryID,
	}

	if err := db.Create(&service).Error; err != nil {
		log.Printf("Erro ao criar serviço no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Serviço criado com sucesso: ID %d", service.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Serviço criado com sucesso",
		"service": service,
	})
}
