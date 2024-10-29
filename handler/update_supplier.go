package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

type UpdateSupplierInput struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	ServiceIDs []uint `json:"service_ids" binding:"required"`
	ProductIDs []uint `json:"product_ids"`
}

// @BasePath /api/v1

// @Summary Update supplier
// @Description Update a supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier Identification"
// @Param supplier body UpdateSupplierInput true "Supplier data to Update"
// @Success 200 {object} UpdateSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers/{id} [put]

// Atualiza um fornecedor existente
func UpdateSupplierHandler(c *gin.Context) {
	supplierID := c.Param("id")
	var input UpdateSupplierInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro no binding do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Atualizando fornecedor %s com categoria %d, serviços %v e produtos %v",
		supplierID, input.CategoryID, input.ServiceIDs, input.ProductIDs)

	// Buscar o fornecedor existente
	var supplierLink schemas.SupplierLink
	if err := db.First(&supplierLink, supplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fornecedor não encontrado"})
		return
	}

	// Atualizar categoria
	supplierLink.CategoryID = input.CategoryID
	if err := db.Save(&supplierLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria"})
		return
	}

	// Remover serviços existentes
	if err := db.Where("supplier_link_id = ?", supplierLink.ID).Delete(&schemas.SupplierService{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover serviços antigos"})
		return
	}

	// Adicionar novos serviços
	for _, serviceID := range input.ServiceIDs {
		supplierService := schemas.SupplierService{
			SupplierLinkID: supplierLink.ID,
			ServiceID:      serviceID,
		}
		if err := db.Create(&supplierService).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar serviço"})
			return
		}
	}

	// Remover produtos existentes
	if err := db.Where("supplier_link_id = ?", supplierLink.ID).Delete(&schemas.SupplierProduct{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover produtos antigos"})
		return
	}

	// Adicionar novos produtos
	for _, productID := range input.ProductIDs {
		supplierProduct := schemas.SupplierProduct{
			SupplierLinkID: supplierLink.ID,
			ProductID:      productID,
		}
		if err := db.Create(&supplierProduct).Error; err != nil {
			log.Printf("Erro ao criar vínculo com produto %d: %v", productID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar produto"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fornecedor atualizado com sucesso",
		"supplier": supplierLink,
	})
}
