package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Update Supplier Product
func UpdateProductHandler(c *gin.Context) {
	log.Println("Iniciando UpdateSupplierProduct")

	// Recebe o ID do produto a ser atualizado
	var input struct {
		ID         uint   `json:"id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		ServiceID uint   `json:"service_id" binding:"required"` 
	}

	// Faz o bind do JSON recebido para a estrutura
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input recebido: %+v", input)

	// Busca o produto no banco de dados
	var product schemas.Product

	if err := db.First(&product, input.ID).Error; err != nil {
		log.Printf("Produto não encontrado: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	// Atualiza os campos do produto
	product.Name = input.Name 
	product.ServiceID = input.ServiceID 

	// Salva as alterações no banco de dados
	if err := db.Save(&product).Error; err != nil {
		log.Printf("Erro ao atualizar produto no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar produto"})
		return
	}

	// Após atualizar o produto, busque a lista atualizada de produtos
	var products []schemas.Product
	if err := db.Preload("Service").Find(&products).Error; err != nil {
		log.Printf("Erro ao buscar produtos no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produtos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Produto atualizado com sucesso",
		"product":  product,
		"products": products, // Envie a lista atualizada de produtos
	})

}
