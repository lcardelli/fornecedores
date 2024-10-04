package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// @BasePath /api/v1

// @Summary List suppliers
// @Description List all suppliers
// @Tags Suppliers
// @Accept json
// @Produce json
// @Success 200 {object} ListSuppliersResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [get]
func ListSupplierHandler(ctx *gin.Context) {
	suppliers := []schemas.Supplier{}

	if err := db.Find(&suppliers).Error; err != nil {
		SendError(ctx, http.StatusInternalServerError, "Error listing suppliers")
		return
	}

	SendSucces(ctx, "list-suppliers", suppliers)
}
