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

	supplier := schemas.Supplier{}

	if err := db.First(&supplier, id).Error; err != nil {
		SendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	// Update the supplier
	if request.Name != "" {
		supplier.Name = request.Name
	}

	if request.CNPJ != "" {
		supplier.CNPJ = request.CNPJ
	}

	if request.Email != "" {
		supplier.Email = request.Email
	}

	if request.Phone != "" {
		supplier.Phone = request.Phone
	}

	if request.Address != "" {
		supplier.Address = request.Address
	}

	if request.CategoryID != 0 {
		supplier.CategoryID = request.CategoryID
	}

	if request.Services != nil {
		supplier.Services = []schemas.SupplierService{}

		for _, service := range request.Services {
			supplier.Services = append(supplier.Services, schemas.SupplierService{
				ServiceID: service.ServiceID,
			})
		}
	}

	if err := db.Save(&supplier).Error; err != nil {
		logger.Errorf("error updating supplier: %v", err.Error())
		SendError(ctx, http.StatusInternalServerError, "error updating supplier")
		return
	}

	SendSucces(ctx, "Update Supplier", supplier)
}
