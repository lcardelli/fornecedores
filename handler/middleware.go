package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID") // Obtém o ID do usuário da sessão
		if userID == nil {
			// Redireciona para a página de login se o usuário não estiver autenticado
			c.Redirect(http.StatusFound, "/api/v1/index?error=unauthorized") // Adiciona um parâmetro de erro à URL
			c.Abort()                                                        // Interrompe a execução da requisição
			return
		}

		var user schemas.User
		if err := db.First(&user, userID).Error; err != nil {
			c.Redirect(http.StatusFound, "/api/v1/index?error=user_not_found")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
