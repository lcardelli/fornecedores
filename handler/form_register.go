package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func FormRegisterHandler(c *gin.Context) {
	user, _ := c.Get("user")
	typedUser := user.(schemas.User)

	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	fornecedores, err := getFornecedoresFromDatabase()
	if err != nil {
		log.Printf("Erro ao buscar fornecedores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados de fornecedores: " + err.Error()})
		return
	}

	log.Printf("Número de fornecedores passados para o template: %d", len(fornecedores))

	c.HTML(http.StatusOK, "form_register.html", gin.H{
		"user":         typedUser,
		"Categories":   categories,
		"Services":     services,
		"Fornecedores": fornecedores,
		"activeMenu":   "cadastro-fornecedor",
	})
}
