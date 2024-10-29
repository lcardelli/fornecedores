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

// Mostra um fornecedor pelo ID
func ShowSupplierHandler(ctx *gin.Context) {
	supplierID := ctx.Query("id")

	// Verifica se o ID foi fornecido
	if supplierID == "" {
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
		return
	}

	// Estrutura para armazenar o fornecedor
	supplierLink := schemas.SupplierLink{}

	// Busca o fornecedor no banco de dados
	if err := db.Preload("Category").Preload("Services").Preload("Services.Service").First(&supplierLink, "id = ?", supplierID).Error; err != nil {
		SendError(ctx, http.StatusNotFound, "Supplier link not found")
		return
	}

	// Busca informações do fornecedor externo
	externalSupplier, err := getFornecedorByCNPJ(supplierLink.CNPJ)
	if err != nil {
		SendError(ctx, http.StatusNotFound, "External supplier information not found")
		return
	}

	// Cria a resposta com as informações do fornecedor
	response := schemas.SupplierLinkResponse{
		ID:               supplierLink.ID,
		CNPJ:             supplierLink.CNPJ,
		CategoryID:       supplierLink.CategoryID,
		Category:         supplierLink.Category,
		Services:         convertToServiceResponses(supplierLink.Services),
		Products:         convertToProductResponses(supplierLink.Products),
		CreatedAt:        supplierLink.CreatedAt,
		UpdatedAt:        supplierLink.UpdatedAt,
		DeletedAt:        supplierLink.DeletedAt.Time,
		ExternalSupplier: *externalSupplier,
	}

	// Envia a resposta de sucesso
	SendSucces(ctx, "show-supplier", response)
}

func convertToProductResponses(products []schemas.SupplierProduct) []schemas.ProductResponse {
	var responses []schemas.ProductResponse
	for _, product := range products {
		responses = append(responses, schemas.ProductResponse{
			ID:        product.Product.ID,
			Name:      product.Product.Name,
		})
	}
	return responses
}
