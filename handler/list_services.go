package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Lista todos os serviços
func ListServicesHandler(c *gin.Context) {
	var services []schemas.Service
	// Busca todos os serviços com a categoria correspondente
	if err := db.Preload("Category").Find(&services).Error; err != nil {
		log.Printf("Erro ao buscar serviços: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	// Criar uma slice de ServiceResponse
	var serviceResponses []schemas.ServiceResponse
	for _, service := range services {
		serviceResponse := schemas.ServiceResponse{
			ID:         service.ID,
			Name:       service.Name,
			CategoryID: service.CategoryID,
			Category: schemas.SupplierCategoryResponse{
				ID:   service.Category.ID,
				Name: service.Category.Name,
			},
		}
		serviceResponses = append(serviceResponses, serviceResponse)
	}

	c.JSON(http.StatusOK, serviceResponses)
}
