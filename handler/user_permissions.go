package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// GetUserPermissionsHandler retorna as permissões do usuário
func GetUserPermissionsHandler(c *gin.Context) {
	// Verificar se o usuário atual é admin global
	currentUser, _ := c.Get("user")
	userModel := currentUser.(schemas.User)
	if !userModel.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Apenas administradores globais podem visualizar permissões"})
		return
	}

	userId := c.Param("id")
	var department schemas.UserDepartment
	var user schemas.User

	// Primeiro busca o usuário para pegar o status de admin
	if err := db.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	result := db.Where("user_id = ?", userId).First(&department)
	if result.Error != nil {
		// Retorna permissões vazias se não encontrar, mas mantém o status de admin
		c.JSON(http.StatusOK, gin.H{
			"is_admin":        user.Admin,
			"department":      "",
			"view_suppliers":  false,
			"view_licenses":   false,
			"admin_suppliers": false,
			"admin_licenses":  false,
		})
		return
	}

	var departmentName string
	if department.DepartmentID > 0 {
		var dept schemas.Departament
		if err := db.First(&dept, department.DepartmentID).Error; err == nil {
			departmentName = dept.Name
		}
	}

	// Retorna todas as permissões incluindo o status de admin
	c.JSON(http.StatusOK, gin.H{
		"is_admin":        user.Admin,
		"department":      departmentName,
		"view_suppliers":  department.ViewSuppliers,
		"view_licenses":   department.ViewLicenses,
		"admin_suppliers": department.AdminSuppliers,
		"admin_licenses":  department.AdminLicenses,
	})
}

// UpdateUserPermissionsHandler atualiza as permissões do usuário
func UpdateUserPermissionsHandler(c *gin.Context) {
	var req struct {
		UserID         uint   `json:"user_id"`
		Department     string `json:"department"`
		ViewSuppliers  bool   `json:"view_suppliers"`
		ViewLicenses   bool   `json:"view_licenses"`
		AdminSuppliers bool   `json:"admin_suppliers"`
		AdminLicenses  bool   `json:"admin_licenses"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Erro ao fazer bind do JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Atualizando permissões - Dados recebidos: %+v\n", req)

	// Buscar registro existente ou criar novo
	var userDepartment schemas.UserDepartment
	result := db.Where("user_id = ?", req.UserID).First(&userDepartment)
	
	if result.Error != nil {
		// Se não encontrar, cria novo registro
		userDepartment = schemas.UserDepartment{
			UserID:       req.UserID,
			DepartmentID: 1, // Departamento padrão
		}
	}

	// Atualiza apenas os campos específicos
	userDepartment.ViewSuppliers = req.ViewSuppliers
	userDepartment.AdminSuppliers = req.AdminSuppliers
	userDepartment.ViewLicenses = req.ViewLicenses
	userDepartment.AdminLicenses = req.AdminLicenses

	// Salva as alterações
	if result.Error != nil {
		if err := db.Create(&userDepartment).Error; err != nil {
			fmt.Printf("Erro ao criar permissões: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := db.Save(&userDepartment).Error; err != nil {
			fmt.Printf("Erro ao atualizar permissões: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	fmt.Printf("Permissões atualizadas com sucesso: %+v\n", userDepartment)
	c.JSON(http.StatusOK, gin.H{
		"message": "Permissões atualizadas com sucesso",
		"data":    userDepartment,
	})
}
