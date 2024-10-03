package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Delete Supplier Handler
func DeleteSupplierHandler(ctx *gin.Context) {
	supplierID := ctx.Query("id") // Obtém o ID do fornecedor a partir dos parâmetros da URL

	// Verify if id is not null
	if supplierID == "" {
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}
	// Verify if existing supplier
	var supplier schemas.Supplier
	if err := db.First(&supplier, supplierID).Error; err != nil {
		logger.Errorf("Error finding supplier: %v", err)
		SendError(ctx, http.StatusNotFound, fmt.Sprintf("supplier with id: %s not found", supplierID))
		return
	}

	// Excluir fornecedor fisicamente (os serviços associados serão excluídos em cascata)
	if err := db.Delete(&supplier).Error; err != nil {
		logger.Errorf("Error deleting supplier: %v", err)
		SendError(ctx, http.StatusInternalServerError, fmt.Sprintf("Error deleting supplier with id: %s", supplierID))
		return
	}

	SendSucces(ctx, "delete-supplier", supplier) // Retorna sucesso com o ID do fornecedor deletado
}
