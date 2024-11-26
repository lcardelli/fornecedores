package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// RenderViewLicensesPage renderiza a página de visualização de licenças
func RenderViewLicensesPage(c *gin.Context) {
	RenderTemplate(c, "list_licenses.html", gin.H{
		"activeMenu": "visualizar-licencas",
	})
}

// ListLicensesHandler processa a requisição de listagem de licenças
func ListLicensesHandler(c *gin.Context) {
	// Obter parâmetros de filtro
	search := c.Query("search")
	status := c.Query("status")
	dateFilter := c.Query("date")

	licenses := GetFilteredLicenses(search, status, dateFilter)

	c.JSON(http.StatusOK, gin.H{
		"licenses": licenses,
	})
}

// GetFilteredLicenses retorna as licenças filtradas de acordo com os parâmetros
func GetFilteredLicenses(search, status, dateFilter string) []schemas.License {
	query := db.Table("licenses").
		Preload("Software") // Carrega os dados do software relacionado

	// Aplicar filtros
	if search != "" {
		query = query.Where(
			"license_key LIKE ? OR software.name LIKE ?",
			"%"+search+"%", "%"+search+"%",
		).Joins("LEFT JOIN software ON licenses.software_id = software.id")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filtro de data
	now := time.Now()
	switch dateFilter {
	case "today":
		query = query.Where("DATE(created_at) = DATE(?)", now)
	case "week":
		query = query.Where("created_at >= ?", now.AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", now.AddDate(0, -1, 0))
	case "year":
		query = query.Where("created_at >= ?", now.AddDate(-1, 0, 0))
	}

	var licenses []schemas.License
	query.Find(&licenses)

	// O status é calculado automaticamente pelos hooks do modelo
	// através do método CalculateStatus

	return licenses
}
