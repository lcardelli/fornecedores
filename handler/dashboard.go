package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/lcardelli/fornecedores/schemas"
)

// DashboardHandler renderiza o template do dashboard com as informações do usuário
func DashboardHandler(c *gin.Context) {
    // Supondo que você tenha um middleware que armazena o ID do usuário na sessão
    userID := c.MustGet("userID").(string) // Ajuste conforme sua implementação de sessão

    var user schemas.User
    if err := db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.HTML(http.StatusOK, "dashboard.html", gin.H{
        "user": user,
    })
}