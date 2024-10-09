package handler

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"html/template"
)

func DashboardHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID") // Obtém o ID do usuário da sessão

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user schemas.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	supplierCountByCategory, err := GetSupplierCountByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados"})
		return
	}

	// Carregar templates
	tmpl := template.Must(template.ParseGlob("templates/*"))
	tmpl.ExecuteTemplate(c.Writer, "dashboard.html", gin.H{
		"user":                    user,
		"supplierCountByCategory": supplierCountByCategory,
	})
}

// GetSupplierCountByCategory retorna a contagem de fornecedores por categoria
func GetSupplierCountByCategory() ([]struct { 	
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
}, error) {
	var results []struct {
		CategoryName string `json:"category_name"`
		Count        int    `json:"count"`
	}

	if err := db.Table("suppliers").
		Select("supplier_categories.name as category_name, COUNT(suppliers.id) as count").
		Joins("LEFT JOIN supplier_categories ON suppliers.category_id = supplier_categories.id").
		Group("supplier_categories.name").
		Scan(&results).Error; err != nil {
		return nil, err 
	}

	return results, nil
}
