package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Cria um novo serviço
func CreateServiceHandler(c *gin.Context) {
	log.Println("Iniciando CreateServiceHandler")

	var service schemas.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&service).Error; err != nil {
		log.Printf("Erro ao criar serviço no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar serviço"})
		return
	}

	log.Printf("Serviço criado com sucesso: ID %d", service.ID)
	c.JSON(http.StatusCreated, service)
}
