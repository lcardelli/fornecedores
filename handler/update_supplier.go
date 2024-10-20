package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// @BasePath /api/v1

// @Summary Update supplier
// @Description Update a supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier Identification"
// @Param supplier body UpdateSupplierRequest true "Supplier data to Update"
// @Success 200 {object} UpdateSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{id} [put]

// Atualiza um fornecedor existente
func UpdateSupplierHandler(ctx *gin.Context) {
	request := UpdateSupplierRequest{}

	// Faz o bind do JSON recebido para a estrutura
	if err := ctx.ShouldBindJSON(&request); err != nil {
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Valida os dados do fornecedor
	if err := request.Validate(); err != nil {
		logger.Errorf("request validation failed: %v", err.Error())
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Obtém o ID do fornecedor a partir do parâmetro da URL
	id := ctx.Param("id")
	// Verifica se o ID foi fornecido
	if id == "" {
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "path").Error())
		return
	}

	// Busca o fornecedor no banco de dados
	supplierLink := schemas.SupplierLink{}

	if err := db.First(&supplierLink, id).Error; err != nil {
		SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	// Atualiza a categoria do fornecedor
	if request.CategoryID != 0 {
		supplierLink.CategoryID = request.CategoryID
	}

	// Atualiza os serviços do fornecedor
	if request.Services != nil {
		// Remover serviços existentes
		if err := db.Where("supplier_link_id = ?", supplierLink.ID).Delete(&schemas.SupplierService{}).Error; err != nil {
			SendError(ctx, http.StatusInternalServerError, "error removing existing services")
			return
		}

		// Adicionar novos serviços
		for _, service := range request.Services {
			newService := schemas.SupplierService{
				SupplierLinkID: supplierLink.ID,
				ServiceID:      service.ServiceID,
			}
			if err := db.Create(&newService).Error; err != nil {
				SendError(ctx, http.StatusInternalServerError, "error adding new services")
				return
			}
		}
	}

	// Salva as alterações no banco de dados
	if err := db.Save(&supplierLink).Error; err != nil {
		logger.Errorf("error updating supplier link: %v", err.Error())
		SendError(ctx, http.StatusInternalServerError, "error updating supplier link")
		return
	}

	// Envia a resposta de sucesso
	SendSucces(ctx, "Update Supplier", supplierLink)
}
