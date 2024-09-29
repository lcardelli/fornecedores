package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func CreateSupplierHandler(ctx *gin.Context) {
	request := CreateSupplierRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		logger.Errorf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("Error validating request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifique se a categoria existe
	var categoryExists bool
	if err := db.Table("supplier_categories").Select("count(*) > 0").Where("id = ?", request.CategoryID).Scan(&categoryExists).Error; err != nil {
		logger.Errorf("Error checking category existence: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking category"})
		return
	}

	if !categoryExists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Category does not exist"})
		return
	}

	// Verifique se o CNPJ já existe
	var cnpjExists bool
	if err := db.Table("suppliers").Select("count(*) > 0").Where("cnpj = ?", request.CNPJ).Scan(&cnpjExists).Error; err != nil {
		logger.Errorf("Error checking CNPJ existence: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking CNPJ"})
		return
	}

	if cnpjExists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CNPJ already exists"})
		return
	}
	
	// Primeiro, crie o fornecedor
	supplier := schemas.Supplier{
		Name:       request.Name,
		CNPJ:       request.CNPJ,
		Email:      request.Email,
		Phone:      request.Phone,
		Address:    request.Address,
		CategoryID: request.CategoryID,
	}

	// Crie o fornecedor no banco de dados
	if err := db.Create(&supplier).Error; err != nil {
		logger.Errorf("Error creating supplier: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating supplier"})
		return
	}

	// Agora, crie os serviços associados
	for i := range request.Services {
		// Atribua o SupplierID do fornecedor recém-criado
		request.Services[i].SupplierID = supplier.ID // O ID do fornecedor é gerado automaticamente

		// Verifique se o serviço já existe
		var serviceExists bool
		if err := db.Table("supplier_services").Select("count(*) > 0").Where("name = ? AND supplier_id = ?", request.Services[i].Name, request.Services[i].SupplierID).Scan(&serviceExists).Error; err != nil {
			logger.Errorf("Error checking service existence: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking service"})
			return
		}

		if serviceExists {
			logger.Warnf("Service %s already exists for supplier %d", request.Services[i].Name, request.Services[i].SupplierID)
			continue // Pula a inserção se o serviço já existir
		}

		// Insira o serviço
		if err := db.Create(&request.Services[i]).Error; err != nil {
			logger.Errorf("Error creating service: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating service"})
			return
		}
	}
}
