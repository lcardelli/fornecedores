package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// GetUserPermissionsHandler retorna as permissões do usuário
func GetUserPermissionsHandler(c *gin.Context) {
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
	
	// Retorna todas as permissões incluindo o status de admin
	c.JSON(http.StatusOK, gin.H{
		"is_admin":        user.Admin,
		"department":      department.Department,
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

	fmt.Printf("Dados recebidos: %+v\n", req)

	var department schemas.UserDepartment
	result := db.Where("user_id = ?", req.UserID).First(&department)
	
	if result.Error != nil {
		// Se não encontrar, cria um novo registro
		department = schemas.UserDepartment{
			UserID:         req.UserID,
			Department:     req.Department,
			ViewSuppliers:  req.ViewSuppliers,
			ViewLicenses:   req.ViewLicenses,
			AdminSuppliers: req.AdminSuppliers,
			AdminLicenses:  req.AdminLicenses,
		}
		if err := db.Create(&department).Error; err != nil {
			fmt.Printf("Erro ao criar permissões: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Atualiza o registro existente
		department.Department = req.Department
		department.ViewSuppliers = req.ViewSuppliers
		department.ViewLicenses = req.ViewLicenses
		department.AdminSuppliers = req.AdminSuppliers
		department.AdminLicenses = req.AdminLicenses
		
		if err := db.Save(&department).Error; err != nil {
			fmt.Printf("Erro ao atualizar permissões: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	fmt.Printf("Permissões atualizadas com sucesso: %+v\n", department)
	c.JSON(http.StatusOK, gin.H{
		"message": "Permissões atualizadas com sucesso",
		"data": department,
	})
}
