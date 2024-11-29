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

		// Adicionar permissões ao contexto
		var userDepartment schemas.UserDepartment
		result := db.Where("user_id = ?", user.ID).First(&userDepartment)
		if result.Error != nil {
			// Se não encontrar permissões, cria um registro vazio
			userDepartment = schemas.UserDepartment{
				UserID: user.ID,
			}
			// Tenta buscar novamente para garantir
			db.Where("user_id = ?", user.ID).FirstOrCreate(&userDepartment)
		}

		fmt.Printf("Permissões carregadas para usuário %d: %+v\n", user.ID, userDepartment)

		c.Set("user", user)
		c.Set("userDepartment", userDepartment)
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

		userDepartment, exists := c.Get("userDepartment")
		if !exists {
			RenderTemplate(c, "permission.html", gin.H{
				"message": "Erro ao carregar permissões",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		department := userDepartment.(schemas.UserDepartment)
		fmt.Printf("Verificando permissão %s para usuário %d\n", permission, userModel.ID)
		fmt.Printf("Permissões do usuário: %+v\n", department)

		// Verifica cada permissão independentemente
		var hasAccess bool
		var message string

		switch permission {
		case "suppliers":
			hasAccess = department.ViewSuppliers || department.AdminSuppliers
			message = "Você não tem permissão para acessar a área de fornecedores"
		case "licenses":
			hasAccess = department.ViewLicenses || department.AdminLicenses
			message = "Você não tem permissão para acessar a área de licenças"
		case "supplier_admin":
			hasAccess = department.AdminSuppliers
			message = "Você não tem permissão para administrar fornecedores"
		case "license_admin":
			hasAccess = department.AdminLicenses
			message = "Você não tem permissão para administrar licenças"
		}

		if !hasAccess {
			fmt.Printf("Acesso negado: %s\n", message)
			RenderTemplate(c, "permission.html", gin.H{
				"message": message,
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		fmt.Printf("Acesso permitido para %s\n", permission)
		c.Next()
	}
}

// SupplierAdminMiddleware verifica APENAS permissões de admin de fornecedores
func SupplierAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userModel := user.(schemas.User)
		
		// Se for admin global, permite acesso
		if userModel.Admin {
			c.Next()
			return
		}

		// Buscar permissões diretamente do banco
		var department schemas.UserDepartment
		result := db.Where("user_id = ?", userModel.ID).First(&department)
		
		// Verifica APENAS AdminSuppliers
		if result.Error != nil || !department.AdminSuppliers {
			RenderTemplate(c, "permission.html", gin.H{
				"message": "Acesso negado: você precisa ser administrador de fornecedores",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LicenseAdminMiddleware verifica APENAS permissões de admin de licenças
func LicenseAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		userModel := user.(schemas.User)
		
		// Se for admin global, permite acesso
		if userModel.Admin {
			c.Next()
			return
		}

		// Buscar permissões diretamente do banco
		var department schemas.UserDepartment
		result := db.Where("user_id = ?", userModel.ID).First(&department)
		
		// Verifica APENAS AdminLicenses
		if result.Error != nil || !department.AdminLicenses {
			RenderTemplate(c, "permission.html", gin.H{
				"message": "Acesso negado: você precisa ser administrador de licenças",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

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
				"message":    "Acesso negado: você precisa ser administrador global do sistema",
				"activeMenu": "dashboard",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
