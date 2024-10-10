package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func ListServicesHandler(c *gin.Context) {
	var services []schemas.Service
	if err := db.Find(&services).Error; err != nil {
		SendError(c, http.StatusInternalServerError, "Error fetching services")
		return
	}
	SendSucces(c, "list-services", services)
}