package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// UpdateServiceHandler lida com a atualização de um serviço existente
func UpdateServiceHandler(c *gin.Context) {
	log.Println("Iniciando UpdateServiceHandler")

	var input struct {
		ID         uint   `json:"id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var service schemas.Service
	if err := db.First(&service, input.ID).Error; err != nil {
		log.Printf("Serviço não encontrado: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	// Atualizar os campos do serviço
	service.Name = input.Name
	service.CategoryID = input.CategoryID

	if err := db.Save(&service).Error; err != nil {
		log.Printf("Erro ao atualizar serviço no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar serviço"})
		return
	}

	// Após atualizar o serviço, busque a lista atualizada de serviços
	var services []schemas.Service
	if err := db.Preload("Category").Find(&services).Error; err != nil {
		log.Printf("Erro ao buscar serviços no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	log.Printf("Serviço atualizado com sucesso: ID %d", service.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Serviço atualizado com sucesso",
		"service":  service,
		"services": services, // Envie a lista atualizada de serviços
	})
}
