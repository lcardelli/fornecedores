package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/gorm"
)



func DeleteCategoryHandler(c *gin.Context) {
	id := c.Param("id")

	// Verifica se o ID é válido
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria não fornecido"})
		return
	}

	var category schemas.SupplierCategory
	result := db.First(&category, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Categoria não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categoria"})
		}
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar categoria"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria deletada com sucesso"})
}

func DeleteMultipleCategories(c *gin.Context) {
	var request struct {
		IDs []int `json:"ids"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	for _, id := range request.IDs {
		if err := db.Delete(&schemas.SupplierCategory{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir categorias"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categorias excluídas com sucesso"})
}
