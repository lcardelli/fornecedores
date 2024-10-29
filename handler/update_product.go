package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Update Supplier Product
func UpdateProductHandler(c *gin.Context) {
	// Recebe o input do usuário
	var input struct {
		Name      string `json:"name" binding:"required"`
		ServiceID uint   `json:"service_id" binding:"required"`
	}

	// Valida o input do usuário
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro no binding do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtém o ID do produto da URL
	productID := c.Param("id")

	// Busca o produto existente
	var product schemas.Product
	if err := db.First(&product, productID).Error; err != nil {
		log.Printf("Erro ao buscar produto: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	// Verifica se o serviço existe
	var service schemas.Service
	if err := db.First(&service, input.ServiceID).Error; err != nil {
		log.Printf("Erro ao buscar serviço: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Serviço não encontrado"})
		return
	}

	// Atualiza os campos do produto
	product.Name = input.Name
	product.ServiceID = input.ServiceID

	// Salva as alterações
	if err := db.Save(&product).Error; err != nil {
		log.Printf("Erro ao atualizar produto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar produto"})
		return
	}

	// Busca o produto atualizado com suas relações
	var updatedProduct schemas.Product
	if err := db.Preload("Service").First(&updatedProduct, product.ID).Error; err != nil {
		log.Printf("Erro ao buscar produto atualizado: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produto atualizado"})
		return
	}

	// Cria a resposta formatada
	response := schemas.ProductResponse{
		ID:        updatedProduct.ID,
		Name:      updatedProduct.Name,
		ServiceID: updatedProduct.ServiceID,
		Service: schemas.ServiceResponse{
			ID:   updatedProduct.Service.ID,
			Name: updatedProduct.Service.Name,
		},
	}

	log.Printf("Produto atualizado com sucesso: ID %d", product.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Produto atualizado com sucesso",
		"product": response,
	})
}
