package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// @BasePath /api/v1

// @Summary Delete Supplier
// @Description Delete a new supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param id query string true "Supplier identification"
// @Success 200 {object} DeleteSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /suppliers [delete]
func DeleteSupplierHandler(ctx *gin.Context) {
	supplierID := ctx.Query("id")

	if supplierID == "" {
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}

	var supplierLink schemas.SupplierLink
	if err := db.First(&supplierLink, supplierID).Error; err != nil {
		logger.Errorf("Error finding supplier link: %v", err)
		SendError(ctx, http.StatusNotFound, fmt.Sprintf("supplier link with id: %s not found", supplierID))
		return
	}

	if err := db.Delete(&supplierLink).Error; err != nil {
		logger.Errorf("Error deleting supplier link: %v", err)
		SendError(ctx, http.StatusInternalServerError, fmt.Sprintf("Error deleting supplier link with id: %s", supplierID))
		return
	}

	SendSucces(ctx, "delete-supplier", supplierLink)
}
