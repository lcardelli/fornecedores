package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func GetServiceListHandler(c *gin.Context) {
	var services []schemas.Service
	if err := db.Preload("Category").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	c.JSON(http.StatusOK, services)
}
