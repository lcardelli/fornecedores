package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// DashboardHandler exibe o dashboard do usuário
func DashboardHandler(c *gin.Context) {
	// Obter o usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do usuário"})
		return
	}

	// Renderizar o template dashboard.html
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user": user,
		// "supplierCountByCategory": supplierCountByCategory,
		"activeMenu": "dashboard",
	})
}

