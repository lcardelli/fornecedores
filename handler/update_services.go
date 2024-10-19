package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Atualiza um serviço existente
func UpdateServiceHandler(c *gin.Context) {
	id := c.Param("id")
	log.Printf("ID recebido para atualização: %s", id)

	// Busca o serviço no banco de dados
	var service schemas.Service
	if err := db.Preload("Category").First(&service, id).Error; err != nil {
		log.Printf("Erro ao buscar serviço: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}
	// Estrutura para receber os dados do serviço
	var updateData struct {
		Name       string `json:"name"`
		CategoryID uint   `json:"category_id"`
	}

	// Faz o bind do JSON recebido para a estrutura
	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza os campos do serviço
	service.Name = updateData.Name
	service.CategoryID = updateData.CategoryID

	// Salva as alterações no banco de dados
	if err := db.Save(&service).Error; err != nil {
		log.Printf("Erro ao atualizar serviço: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar serviço"})
		return
	}

	// Recarrega o serviço com a categoria atualizada
	if err := db.Preload("Category").First(&service, id).Error; err != nil {
		log.Printf("Erro ao recarregar serviço atualizado: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recarregar serviço atualizado"})
		return
	}
	// Log do serviço atualizado
	log.Printf("Serviço atualizado com sucesso: %+v", service)
	c.JSON(http.StatusOK, service)
}