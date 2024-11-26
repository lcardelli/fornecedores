package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// GetUserPermissionsHandler retorna as permissões do usuário
func GetUserPermissionsHandler(c *gin.Context) {
	userId := c.Param("id")
	var department schemas.UserDepartment

	result := db.Where("user_id = ?", userId).First(&department)
	if result.Error != nil {
		// Retorna permissões vazias se não encontrar
		c.JSON(http.StatusOK, gin.H{
			"department":     "",
			"view_suppliers": false,
			"view_licenses":  false,
		})
		return
	}

	c.JSON(http.StatusOK, department)
}

// UpdateUserPermissionsHandler atualiza as permissões do usuário
func UpdateUserPermissionsHandler(c *gin.Context) {
	var req struct {
		UserID        uint   `json:"user_id"`
		Department    string `json:"department"`
		ViewSuppliers bool   `json:"view_suppliers"`
		ViewLicenses  bool   `json:"view_licenses"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user schemas.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	if err := SetDepartmentAccess(&user, req.Department, req.ViewSuppliers, req.ViewLicenses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissões atualizadas com sucesso"})
}
