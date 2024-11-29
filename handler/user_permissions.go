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
	// Inicializa departamentos se necessário
	if err := InitializeDepartments(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inicializar departamentos"})
		return
	}

	// Verificar se o usuário atual é admin global
	currentUser, _ := c.Get("user")
	userModel := currentUser.(schemas.User)
	if !userModel.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Apenas administradores globais podem modificar permissões"})
		return
	}

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

	deptID, err := GetDepartmentIDByName(req.Department)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Departamento inválido"})
		return
	}

	var department schemas.UserDepartment
	result := db.Where("user_id = ?", req.UserID).First(&department)

	if result.Error != nil {
		// Se não encontrar, cria um novo registro
		department = schemas.UserDepartment{
			UserID:         req.UserID,
			DepartmentID:   deptID,
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
		department.DepartmentID = deptID
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

	var checkDepartment schemas.UserDepartment
	db.Where("user_id = ?", req.UserID).First(&checkDepartment)
	fmt.Printf("Permissões após atualização: %+v\n", checkDepartment)

	fmt.Printf("Permissões atualizadas com sucesso: %+v\n", department)
	c.JSON(http.StatusOK, gin.H{
		"message": "Permissões atualizadas com sucesso",
		"data":    department,
	})
}
