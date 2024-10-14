package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func FormRegisterHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID") // Obtém o ID do usuário da sessão

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user schemas.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

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

	c.HTML(http.StatusOK, "form_register.html", gin.H{ 
		"user":       user,
		"Categories": categories,
		"Services":   services,
	})
}
