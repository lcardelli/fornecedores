package handler

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"html/template"
)

// DashboardHandler renderiza o template do dashboard com as informações do usuário
func DashboardHandler(c *gin.Context) {
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

	// Carregar templates
	tmpl := template.Must(template.ParseGlob("templates/*"))
	tmpl.ExecuteTemplate(c.Writer, "dashboard.html", gin.H{
		"user": user,
	})
}
