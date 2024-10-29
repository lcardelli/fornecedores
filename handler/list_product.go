package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// List Supplier Products
func ListSupplierProducts(c *gin.Context) {
	log.Println("Iniciando ListSupplierProducts")

	var products []schemas.Product
	if err := db.Preload("Service").Find(&products).Error; err != nil {
		log.Printf("Erro ao buscar produtos no banco de dados: %v", err)
	}

	var productResponses []schemas.ProductResponse
	for _, product := range products {
		productResponse := schemas.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			ServiceID:   product.ServiceID,
			Service: schemas.ServiceResponse{
				ID:   product.Service.ID,
				Name: product.Service.Name,
			},
		}
		productResponses = append(productResponses, productResponse)
	}

	// Retorna a lista de produtos
	c.JSON(http.StatusOK, gin.H{
		"message": "Produtos listados com sucesso",
		"products": productResponses,
	})
}

