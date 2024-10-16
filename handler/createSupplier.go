package handler

import (
	"fmt"
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
func CreateSupplierHandler(c *gin.Context) {
	var input struct {
		CNPJ       string `form:"cnpj" binding:"required"`
		CategoryID uint   `form:"category_id" binding:"required"`
		ServiceIDs []uint `form:"service_ids[]"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se o fornecedor já está vinculado
	var existingLink schemas.SupplierLink
	if err := db.Where("cnpj = ?", input.CNPJ).First(&existingLink).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fornecedor já vinculado"})
		return
	}

	// Buscar o fornecedor no banco de dados externo
	fornecedor, err := getFornecedorByCNPJ(input.CNPJ)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fornecedor não encontrado no banco de dados externo"})
		return
	}

	// Criar o vínculo do fornecedor
	supplierLink := schemas.SupplierLink{
		CNPJ:       input.CNPJ,
		CategoryID: input.CategoryID,
	}

	if err := db.Create(&supplierLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar vínculo do fornecedor"})
		return
	}

	// Vincular serviços
	for _, serviceID := range input.ServiceIDs {
		supplierService := schemas.SupplierService{
			SupplierLinkID: supplierLink.ID,
			ServiceID:      serviceID,
		}
		if err := db.Create(&supplierService).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao vincular serviço"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Fornecedor vinculado com sucesso", "supplier": fornecedor})
}

func getFornecedorByCNPJ(cnpj string) (*schemas.ExternalSupplier, error) {
	fornecedores, err := getFornecedoresFromDatabase()
	if err != nil {
		return nil, err
	}

	for _, f := range fornecedores {
		if f.CGCCFO.String == cnpj {
			return &schemas.ExternalSupplier{
				CODCOLIGADA:  f.CODCOLIGADA,
				CODCFO:       f.CODCFO,
				NOMEFANTASIA: f.NOMEFANTASIA.String,
				NOME:         f.NOME.String,
				CGCCFO:       f.CGCCFO.String,
				RUA:          f.RUA.String,
				NUMERO:       f.NUMERO.String,
				COMPLEMENTO:  f.COMPLEMENTO.String,
				BAIRRO:       f.BAIRRO.String,
				CIDADE:       f.CIDADE.String,
				CEP:          f.CEP.String,
				TELEFONE:     f.TELEFONE.String,
				EMAIL:        f.EMAIL.String,
				CONTATO:      f.CONTATO.String,
				UF:           f.UF.String,
				ATIVO:        f.ATIVO.String,
				TIPO:         f.TIPO.String,
			}, nil
		}
	}

	return nil, fmt.Errorf("Fornecedor não encontrado")
}
