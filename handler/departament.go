package handler

import (
	"github.com/lcardelli/fornecedores/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitializeDepartments cria os departamentos padrão se não existirem
func InitializeDepartments() error {
	departments := []string{"TI", "Compras", "Geral"}

	for _, deptName := range departments {
		var dept schemas.Departament
		result := db.Where("name = ?", deptName).FirstOrCreate(&dept, schemas.Departament{Name: deptName})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// GetDepartmentIDByName retorna o ID do departamento pelo nome
func GetDepartmentIDByName(name string) (uint, error) {
	var dept schemas.Departament
	result := db.Where("name = ?", name).First(&dept)
	if result.Error != nil {
		return 0, result.Error
	}
	return dept.ID, nil
}

// GetDepartmentsHandler retorna a lista de departamentos
func GetDepartmentsHandler(c *gin.Context) {
	var departments []schemas.Departament
	if err := db.Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, departments)
}
