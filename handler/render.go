package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// GetUserPermissions retorna as permissões do usuário atual
func GetUserPermissions(c *gin.Context) schemas.UserDepartment {
	var department schemas.UserDepartment
	user, _ := c.Get("user")
	userModel := user.(schemas.User)

	db.Where("user_id = ?", userModel.ID).First(&department)
	return department
}

// RenderTemplate renderiza o template com os dados comuns
func RenderTemplate(c *gin.Context, template string, data gin.H) {
	user, _ := c.Get("user")
	userPermissions, _ := c.Get("userPermissions")

	fmt.Printf("RenderTemplate - Template: %s\n", template)
	fmt.Printf("RenderTemplate - User: %+v\n", user)
	fmt.Printf("RenderTemplate - Permissions: %+v\n", userPermissions)

	if data == nil {
		data = gin.H{}
	}

	data["user"] = user
	data["userPermissions"] = userPermissions

	c.HTML(http.StatusOK, template, data)
}

// RenderError renderiza uma página de erro
func RenderError(c *gin.Context, status int, message string) {
	RenderTemplate(c, "error.html", gin.H{
		"status":  status,
		"message": message,
	})
}
