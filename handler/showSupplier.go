package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Show Supplier Handler
func ShowSupplierHandler(ctx *gin.Context) {
	supplierID := ctx.Query("id")

	if supplierID == "" {
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}

	supplier := schemas.Supplier{}

	if err := db.First(&supplier, "id = ?", supplierID).Error; err != nil {
		SendError(ctx, http.StatusNotFound, "Supplier not found")
		return
	}

	SendSucces(ctx, "show-supplier", supplier)
}
