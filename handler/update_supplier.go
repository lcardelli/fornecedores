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
	// Estrutura para receber os dados da requisição
	type UpdateSupplierRequest struct {
		CategoryID uint   `json:"category_id"`
		ServiceIDs []uint `json:"service_ids"`
	}

	var request UpdateSupplierRequest

	// Faz o bind do JSON recebido para a estrutura
	if err := ctx.ShouldBindJSON(&request); err != nil {
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
		SendError(ctx, http.StatusNotFound, "Fornecedor não encontrado")
		return
	}

	// Atualiza a categoria do fornecedor
	supplierLink.CategoryID = request.CategoryID

	// Inicia uma transação
	tx := db.Begin()

	// Atualiza os serviços do fornecedor
	if err := tx.Where("supplier_link_id = ?", supplierLink.ID).Delete(&schemas.SupplierService{}).Error; err != nil {
		tx.Rollback()
		SendError(ctx, http.StatusInternalServerError, "Erro ao remover serviços existentes")
		return
	}

	for _, serviceID := range request.ServiceIDs {
		newService := schemas.SupplierService{
			SupplierLinkID: supplierLink.ID,
			ServiceID:      serviceID,
		}
		if err := tx.Create(&newService).Error; err != nil {
			tx.Rollback()
			SendError(ctx, http.StatusInternalServerError, "Erro ao adicionar novos serviços")
			return
		}
	}

	// Salva as alterações no banco de dados
	if err := tx.Save(&supplierLink).Error; err != nil {
		tx.Rollback()
		logger.Errorf("erro ao atualizar vínculo do fornecedor: %v", err.Error())
		SendError(ctx, http.StatusInternalServerError, "Erro ao atualizar vínculo do fornecedor")
		return
	}

	// Commit da transação
	tx.Commit()

	// Envia a resposta de sucesso
	SendSucces(ctx, "Fornecedor atualizado com sucesso", supplierLink)
}
