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
	var supplierLinks []schemas.SupplierLink

	if err := db.Preload("Category").Preload("Services").Preload("Services.Service").Find(&supplierLinks).Error; err != nil {
		SendError(ctx, http.StatusInternalServerError, "Error listing supplier links")
		return
	}

	var response []schemas.SupplierLinkResponse
	for _, link := range supplierLinks {
		externalSupplier, err := getFornecedorByCNPJ(link.CNPJ)
		if err != nil {
			continue
		}
		response = append(response, schemas.SupplierLinkResponse{
			ID:         link.ID,
			CNPJ:       link.CNPJ,
			CategoryID: link.CategoryID,
			Category:   link.Category,
			Services:   convertToServiceResponses(link.Services),
			CreatedAt:  link.CreatedAt,
			UpdatedAt:  link.UpdatedAt,
			DeletedAt:  link.DeletedAt.Time,
			ExternalSupplier: *externalSupplier,
		})
	}

	SendSucces(ctx, "list-suppliers", response)
}

func convertToServiceResponses(services []schemas.SupplierService) []schemas.ServiceResponse {
	var responses []schemas.ServiceResponse
	for _, service := range services {
		responses = append(responses, schemas.ServiceResponse{
			ID:          service.Service.ID,
			Name:        service.Service.Name,
			Description: service.Service.Description,
		})
	}
	return responses
}
