package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func CreateProductHandler(c *gin.Context) {
	// Recebe o input do usuário
	var input struct {
		Name      string `json:"name" binding:"required"`
		ServiceID uint   `json:"service_id" binding:"required"`
	}

	// Valida o input do usuário
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Verifica se o serviço existe
	var service schemas.Service
	if err := db.Where("id = ?", input.ServiceID).First(&service).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}
	// Cria o produto
	product := schemas.Product{
		Name:      input.Name,
		ServiceID: input.ServiceID,
	}

	// Salva o produto no banco de dados
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})

	// Busca todos os produtos
	var products []schemas.Product
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Envia a lista de produtos atualizada
	log.Printf("Produto criado com sucesso: ID %d", product.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Produto criado com sucesso",
		"product":  product,
		"products": products, // Envie a lista atualizada de produtos
	})

}
