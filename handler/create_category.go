package handler

import (
	"net/http"

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

	category := schemas.SupplierCategory{
		Name: input.Name,
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar categoria"})
		return
	}

	// Após criar a categoria, atualize a lista de categorias no frontend
	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categorias"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Categoria criada com sucesso",
		"category":   category,
		"categories": categories, // Envie a lista atualizada de categorias
	})
}
