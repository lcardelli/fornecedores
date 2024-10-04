package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// @BasePath /api/v1

// @Summary Show supplier
// @Description Show a supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param id query string true "Supplier identification"
// @Success 200 {object} ShowSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /suppliers [get]
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
