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
// @Param id query string true "Supplier Identification"
// @Param supplier body UpdateSupplierRequest true "Supplier data to Update"
// @Success 200 {object} UpdateSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [put]
func UpdateSupplierHandler(ctx *gin.Context) {
	request := UpdateSupplierRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("request validation failed: %v", err.Error())
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id := ctx.Query("id")

	if id == "" {
		SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "string").Error())
		return
	}

	supplierLink := schemas.SupplierLink{}

	if err := db.First(&supplierLink, id).Error; err != nil {
		SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	if request.CategoryID != 0 {
		supplierLink.CategoryID = request.CategoryID
	}

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

	if err := db.Save(&supplierLink).Error; err != nil {
		logger.Errorf("error updating supplier link: %v", err.Error())
		SendError(ctx, http.StatusInternalServerError, "error updating supplier link")
		return
	}

	SendSucces(ctx, "Update Supplier", supplierLink)
}
