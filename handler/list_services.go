package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// ListServicesHandler retorna a lista de todos os serviços
func ListServicesHandler(c *gin.Context) {
	var services []schemas.Service
	
	// Carrega os serviços incluindo a relação com Category
	if err := db.Preload("Category").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	c.JSON(http.StatusOK, services)
}
