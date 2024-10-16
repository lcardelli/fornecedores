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

	supplierLink := schemas.SupplierLink{}

	if err := db.Preload("Category").Preload("Services").Preload("Services.Service").First(&supplierLink, "id = ?", supplierID).Error; err != nil {
		SendError(ctx, http.StatusNotFound, "Supplier link not found")
		return
	}

	externalSupplier, err := getFornecedorByCNPJ(supplierLink.CNPJ)
	if err != nil {
		SendError(ctx, http.StatusNotFound, "External supplier information not found")
		return
	}

	response := schemas.SupplierLinkResponse{
		ID:         supplierLink.ID,
		CNPJ:       supplierLink.CNPJ,
		CategoryID: supplierLink.CategoryID,
		Category:   supplierLink.Category,
		Services:   convertToServiceResponses(supplierLink.Services),
		CreatedAt:  supplierLink.CreatedAt,
		UpdatedAt:  supplierLink.UpdatedAt,
		DeletedAt:  supplierLink.DeletedAt.Time,
			ExternalSupplier: *externalSupplier,
	}

	SendSucces(ctx, "show-supplier", response)
}
