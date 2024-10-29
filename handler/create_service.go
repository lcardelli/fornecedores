package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// CreateServiceHandler lida com a criação de um novo serviço
func CreateServiceHandler(c *gin.Context) {
	log.Println("Iniciando CreateServiceHandler")

	var input struct {
		Name       string `json:"name" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se a categoria existe
	var category schemas.SupplierCategory
	if err := db.First(&category, input.CategoryID).Error; err != nil {
		log.Printf("Categoria não encontrada: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Categoria não encontrada"})
		return
	}

	// Cria o serviço
	service := schemas.Service{
		Name:       input.Name,
		CategoryID: input.CategoryID,
	}

	if err := db.Create(&service).Error; err != nil {
		log.Printf("Erro ao criar serviço no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar serviço"})
		return
	}

	// Após criar o serviço, atualize a lista de serviços no frontend
	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		log.Printf("Erro ao buscar serviços no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	log.Printf("Serviço criado com sucesso: ID %d", service.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Serviço criado com sucesso",
		"service":  service,
		"services": services, // Envie a lista atualizada de serviços
	})
}
