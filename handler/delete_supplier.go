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
	supplierID := ctx.Param("id")
	logger.Infof("Received DELETE request for supplier ID: %s", supplierID)

	if supplierID == "" {
		logger.Error("Supplier ID is empty")
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "path").Error())
		return
	}

	var supplierLink schemas.SupplierLink
	if err := db.First(&supplierLink, supplierID).Error; err != nil {
		logger.Errorf("Error finding supplier link: %v", err)
		SendError(ctx, http.StatusNotFound, fmt.Sprintf("supplier link with id: %s not found", supplierID))
		return
	}

	logger.Infof("Found supplier link to delete: %+v", supplierLink)

	// Use o soft delete
	if err := db.Delete(&supplierLink).Error; err != nil {
		SendError(ctx, http.StatusInternalServerError, "Erro ao deletar fornecedor")
		return
	}

	// Também faça soft delete dos serviços associados
	if err := db.Where("supplier_link_id = ?", supplierID).Delete(&schemas.SupplierService{}).Error; err != nil {
		SendError(ctx, http.StatusInternalServerError, "Erro ao deletar serviços do fornecedor")
		return
	}

	SendSucces(ctx, "Fornecedor deletado com sucesso", nil)
}
