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
	userDepartment, _ := c.Get("userDepartment")

	fmt.Printf("\n=== Debug RenderTemplate ===\n")
	fmt.Printf("Template: %s\n", template)
	fmt.Printf("User: %+v\n", user)
	fmt.Printf("UserDepartment: %+v\n", userDepartment)
	
	// Verificar permissões específicas
	if dept, ok := userDepartment.(schemas.UserDepartment); ok {
		fmt.Printf("Permissões detalhadas:\n")
		fmt.Printf("- ViewSuppliers: %v\n", dept.ViewSuppliers)
		fmt.Printf("- AdminSuppliers: %v\n", dept.AdminSuppliers)
		fmt.Printf("- ViewLicenses: %v\n", dept.ViewLicenses)
		fmt.Printf("- AdminLicenses: %v\n", dept.AdminLicenses)
	} else {
		fmt.Printf("Erro ao converter userDepartment para o tipo correto\n")
	}

	if data == nil {
		data = gin.H{}
	}

	data["user"] = user
	data["userDepartment"] = userDepartment

	fmt.Printf("Data final passada para o template: %+v\n", data)
	fmt.Printf("=== Fim Debug ===\n\n")

	c.HTML(http.StatusOK, template, data)
}

// RenderError renderiza uma página de erro
func RenderError(c *gin.Context, status int, message string) {
	RenderTemplate(c, "error.html", gin.H{
		"status":  status,
		"message": message,
	})
}
