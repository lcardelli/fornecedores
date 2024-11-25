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
		userID := session.Get("userID")
		if userID == nil {
			c.Redirect(http.StatusFound, "/api/v1/index?error=unauthorized")
			c.Abort()
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

// AdminMiddleware verifica se o usuário é administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		user, ok := userInterface.(schemas.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados do usuário"})
			c.Abort()
			return
		}

		if !user.Admin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado. Apenas administradores podem acessar esta área."})
			c.Abort()
			return
		}

		c.Next()
	}
}
