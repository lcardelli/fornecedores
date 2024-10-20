package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler handles GET requests for the home page
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Página Inicial", // Passa dados para o template
	})
	
}
