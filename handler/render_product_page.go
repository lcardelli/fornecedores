package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func RenderProductPageHandler(c *gin.Context) {
	user, _ := c.Get("user")
	typedUser := user.(schemas.User)

	// Buscar todos os serviços para o select
	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	c.HTML(http.StatusOK, "cadastro_produto.html", gin.H{
		"user":       typedUser,
		"Services":   services,
		"activeMenu": "cadastro-produto",
	})
}
