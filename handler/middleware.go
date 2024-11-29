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

		// Se for admin global, permite acesso a tudo
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
				"message": "Você não tem permissão para acessar esta área",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SupplierAdminMiddleware verifica se o usuário é administrador de fornecedores
func SupplierAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userModel := user.(schemas.User)
		fmt.Printf("Verificando permissões de admin fornecedores para usuário %d\n", userModel.ID)
		
		// Se for admin global, permite acesso
		if userModel.Admin {
			c.Next()
			return
		}

		if !HasSupplierAdminAccess(&userModel) {
			fmt.Printf("Acesso negado para admin fornecedores - usuário %d\n", userModel.ID)
			RenderTemplate(c, "permission.html", gin.H{
				"message": "Acesso negado: você precisa ser administrador de fornecedores",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		fmt.Printf("Acesso permitido para admin fornecedores - usuário %d\n", userModel.ID)
		c.Next()
	}
}

// LicenseAdminMiddleware verifica se o usuário é administrador de licenças
func LicenseAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userModel := user.(schemas.User)
		fmt.Printf("Verificando permissões de admin licenças para usuário %d\n", userModel.ID)
		
		// Se for admin global, permite acesso
		if userModel.Admin {
			c.Next()
			return
		}

		if !HasLicenseAdminAccess(&userModel) {
			fmt.Printf("Acesso negado para admin licenças - usuário %d\n", userModel.ID)
			RenderTemplate(c, "permission.html", gin.H{
				"message": "Acesso negado: você precisa ser administrador de licenças",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		fmt.Printf("Acesso permitido para admin licenças - usuário %d\n", userModel.ID)
		c.Next()
	}
}

// GlobalAdminMiddleware verifica se o usuário é administrador global
func GlobalAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userModel := user.(schemas.User)
		fmt.Printf("Verificando admin global para usuário %d: %v\n", userModel.ID, userModel.Admin)
		
		if !userModel.Admin {
			RenderTemplate(c, "permission.html", gin.H{
				"message": "Acesso negado: você precisa ser administrador global do sistema",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
