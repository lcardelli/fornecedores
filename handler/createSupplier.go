package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// @BasePath /api/v1

// @Summary Create Supplier
// @Description Create a new supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param request body CreateSupplierRequest true "Request body"
// @Success 200 {object} CreateSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [post]
func CreateSupplierHandler(ctx *gin.Context) {
	request := CreateSupplierRequest{}

	// Bind the request
	if err := ctx.BindJSON(&request); err != nil {
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the category exists
	var categoryExists bool
	if err := db.Table("supplier_categories").Select("count(*) > 0").Where("id = ?", request.CategoryID).Scan(&categoryExists).Error; err != nil {
		logger.Errorf("Error checking category existence: %v", err)
		SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// If the category does not exist, return an error
	if !categoryExists {
		logger.Warnf("Category %d does not exist", request.CategoryID)
		SendError(ctx, http.StatusBadRequest, "Category does not exist")
		return
	}

	// Check if the CNPJ already exists
	var cnpjExists bool
	if err := db.Table("suppliers").Select("count(*) > 0").Where("cnpj = ?", request.CNPJ).Scan(&cnpjExists).Error; err != nil {
		logger.Errorf("Error checking CNPJ existence: %v", err)
		SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if cnpjExists {
		logger.Warnf("CNPJ %s already exists", request.CNPJ)
		SendError(ctx, http.StatusBadRequest, "CNPJ already exists")
		return
	}

	// First, create the supplier
	supplier := schemas.Supplier{
		Name:       request.Name,
		CNPJ:       request.CNPJ,
		Email:      request.Email,
		Phone:      request.Phone,
		Address:    request.Address,
		CategoryID: request.CategoryID,
	}

	// Create the supplier in the database
	if err := db.Create(&supplier).Error; err != nil {
		logger.Errorf("Error creating supplier: %v", err)
		SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Now, create the associated services
	for _, serviceID := range request.ServiceIDs {
		supplierService := schemas.SupplierService{
			SupplierID: supplier.ID,
			ServiceID: serviceID,
		}
		if err := db.Create(&supplierService).Error; err != nil {
			SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Carregar a categoria e os serviços associados
	if err := db.Preload("Category").Preload("Services").Preload("Services.Service").First(&supplier, supplier.ID).Error; err != nil {
		logger.Errorf("Error loading supplier data: %v", err)
		SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Criar a resposta
	response := schemas.SupplierResponse{
		ID:         supplier.ID,
		Name:       supplier.Name,
		CNPJ:       supplier.CNPJ,
		Email:      supplier.Email,
		Phone:      supplier.Phone,
		Address:    supplier.Address,
		CategoryID: supplier.CategoryID,
		Category:   supplier.Category,
		Services:   make([]schemas.ServiceResponse, len(supplier.Services)),
		CreatedAt:  supplier.CreatedAt,
		UpdatedAt:  supplier.UpdatedAt,
		DeletedAt:  supplier.DeletedAt.Time,
	}

	for i, service := range supplier.Services {
		response.Services[i] = schemas.ServiceResponse{
			ID:          service.Service.ID,
			Name:        service.Service.Name,
			Description: service.Service.Description,
			Price:       service.Service.Price,
		}
	}

	SendSucces(ctx, "create-supplier", response)
}
