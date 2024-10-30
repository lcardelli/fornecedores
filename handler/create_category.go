package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// CreateCategoryHandler lida com a criação de uma nova categoria
func CreateCategoryHandler(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove espaços do início e fim
	input.Name = strings.TrimSpace(input.Name)

	// Verifica se está vazio após o trim
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name cannot be empty"})
		return
	}

	// Verifica se a categoria já existe
	if err := db.Where("name = ?", input.Name).First(&schemas.SupplierCategory{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category already exists"})
		return
	}

	category := schemas.SupplierCategory{
		Name: input.Name,
	}

	if err := db.Create(&category).Error; err != nil {
		log.Printf("Erro ao criar categoria: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Categoria criada com sucesso",
		"category": category,
	})
}
