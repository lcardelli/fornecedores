package handler

import (
	"fmt"
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

		// Carrega as permissões do usuário
		var department schemas.UserDepartment
		result := db.Where("user_id = ?", user.ID).First(&department)
		if result.Error != nil {
			// Se não encontrar permissões, cria um registro vazio
			department = schemas.UserDepartment{
				UserID: user.ID,
			}
		}

		fmt.Printf("Permissões carregadas para usuário %d: %+v\n", user.ID, department)

		c.Set("user", user)
		c.Set("userPermissions", department)
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

// PermissionMiddleware verifica as permissões do usuário
func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		userModel := user.(schemas.User)

		if userModel.Admin {
			c.Next()
			return
		}

		var department schemas.UserDepartment
		db.Where("user_id = ?", userModel.ID).First(&department)

		hasAccess := false
		switch permission {
		case "suppliers":
			hasAccess = department.ViewSuppliers
		case "licenses":
			hasAccess = department.ViewLicenses
		}

		if !hasAccess {
			RenderTemplate(c, "permission.html", gin.H{
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
