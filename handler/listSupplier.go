package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// List Supplier Handler
func ListSupplierHandler(ctx *gin.Context) {
	suppliers := []schemas.Supplier{}

	if err := db.Find(&suppliers).Error; err != nil {
		SendError(ctx, http.StatusInternalServerError, "Error listing suppliers")
		return
	}

	SendSucces(ctx, "list-suppliers", suppliers)
}
