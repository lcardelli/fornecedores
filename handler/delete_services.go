package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/gorm"
)

// Deleta um serviço pelo ID
func DeleteServiceHandler(c *gin.Context) {
	id := c.Param("id")
	log.Printf("ID recebido para deleção: %s", id)

	// Verifica se o ID foi fornecido
	if id == "" || id == "undefined" {
		log.Println("ID do serviço não fornecido ou inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço não fornecido ou inválido"})
		return
	}

	// Converte o ID para uint
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("Erro ao converter ID para uint: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	log.Printf("ID convertido para uint: %d", serviceID)

	// Busca o serviço no banco de dados
	var service schemas.Service
	result := db.First(&service, uint(serviceID))
	// Log do resultado da busca
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Serviço não encontrado para o ID: %d", serviceID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		} else {
			log.Printf("Erro ao buscar serviço: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviço"})
		}
		return
	}

	log.Printf("Serviço encontrado: %+v", service)

	// Deleta o serviço do banco de dados
	if err := db.Delete(&service).Error; err != nil {
		log.Printf("Erro ao deletar serviço: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar serviço"})
		return
	}

	log.Printf("Serviço deletado com sucesso: ID %d", serviceID)
	c.JSON(http.StatusOK, gin.H{"message": "Serviço deletado com sucesso"})
}

// Adicione esta nova função ao seu handler
func DeleteMultipleServices(c *gin.Context) {
	var request struct {
		IDs []int `json:"ids"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	for _, id := range request.IDs {
		if err := db.Delete(&schemas.Service{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir serviços"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Serviços excluídos com sucesso"})
}