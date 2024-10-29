package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func DeleteProductHandler(c *gin.Context) {
	log.Println("Iniciando DeleteProduct")

	id := c.Param("id")
	log.Printf("ID recebido para deleção: %s", id)

	// Verifica se o ID foi fornecido
	if id == "" || id == "undefined" {
		log.Println("ID do produto não fornecido ou inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto não fornecido ou inválido"})
		return
	}

	// Converte o ID para uint
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("Erro ao converter ID para uint: %v", err)
	}
	log.Printf("ID convertido para uint: %d", productID)

	// Busca o produto no banco de dados
	var product schemas.Product
	if err := db.First(&product, productID).Error; err != nil {
		log.Printf("Erro ao buscar produto no banco de dados: %v", err)
	}

	// Deleta o produto
	if err := db.Delete(&product).Error; err != nil {
		log.Printf("Erro ao deletar produto: %v", err)
	}

	log.Println("Produto deletado com sucesso")
	c.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso"})
}

// Deleta vários produtos
func DeleteMultipleProducts(c *gin.Context) {
	var request struct {
		IDs []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	for _, id := range request.IDs {
		if err := db.Delete(&schemas.Product{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir produtos"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produtos excluídos com sucesso"})
}

