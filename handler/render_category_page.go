package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Renderiza a página de cadastro de categorias
func RenderCategoriaHandler(c *gin.Context) {
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

	c.HTML(http.StatusOK, "cadastro_categoria.html", gin.H{
		"user":       user,
		"activeMenu": "cadastro-categoria",
	})
}
