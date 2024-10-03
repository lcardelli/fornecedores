package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// Create Supplier Handler
func CreateSupplierHandler(ctx *gin.Context) {
	request := CreateSupplierRequest{}

	// Bind the request
	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("Error binding JSON: %v", err)
		SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logger.Errorf("Error validating request: %v", err)
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
	for i := range request.Services {
		// Assign the SupplierID of the recently created supplier
		request.Services[i].SupplierID = supplier.ID // O ID do fornecedor é gerado automaticamente

		// Check if the service already exists
		var serviceExists bool
		if err := db.Table("supplier_services").Select("count(*) > 0").Where("name = ? AND supplier_id = ?", request.Services[i].Name, request.Services[i].SupplierID).Scan(&serviceExists).Error; err != nil {
			logger.Errorf("Error checking service existence: %v", err)
			SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		// If the service already exists, skip the insertion
		if serviceExists {
			logger.Warnf("Service %s already exists for supplier %d", request.Services[i].Name, request.Services[i].SupplierID)
			SendError(ctx, http.StatusBadRequest, "Service already exists")
			continue // Pula a inserção se o serviço já existir
		}

		// Insert the service
		if err := db.Create(&request.Services[i]).Error; err != nil {
			logger.Errorf("Error creating service: %v", err)
			SendError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
