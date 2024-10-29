package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/gorm"
)

type CreateSupplierInput struct {
	SupplierCNPJ string   `json:"supplier_cnpj" binding:"required"`
	CategoryID   uint     `json:"category_id" binding:"required"`
	ServiceIDs   []uint   `json:"service_ids" binding:"required"`
	ProductIDs   []uint   `json:"product_ids" binding:"required"`
}

// @BasePath /api/v1

// @Summary Create Supplier
// @Description Create a new supplier
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param request body CreateSupplierInput true "Request body"
// @Success 200 {object} CreateSupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /suppliers [post]

// Cria um novo fornecedor
func CreateSupplierHandler(c *gin.Context) {
	log.Println("Iniciando CreateSupplierHandler")
	var input CreateSupplierInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro ao fazer bind dos dados de entrada: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	log.Printf("Dados recebidos: CNPJ=%s, CategoryID=%d, ServiceIDs=%v, ProductIDs=%v", 
		input.SupplierCNPJ, input.CategoryID, input.ServiceIDs, input.ProductIDs)

	// Verifica se o fornecedor já existe (mesmo que deletado)
	var existingSupplier schemas.SupplierLink
	if err := db.Unscoped().Where("cnpj = ?", input.SupplierCNPJ).First(&existingSupplier).Error; err == nil {
		// Reativa o fornecedor se estiver deletado
		existingSupplier.DeletedAt = gorm.DeletedAt{}
		existingSupplier.CategoryID = input.CategoryID
		
		if err := db.Save(&existingSupplier).Error; err != nil {
			log.Printf("Erro ao reativar fornecedor existente: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao reativar fornecedor"})
			return
		}

		// Remove serviços e produtos existentes
		if err := db.Where("supplier_link_id = ?", existingSupplier.ID).Delete(&schemas.SupplierService{}).Error; err != nil {
			log.Printf("Erro ao remover serviços existentes: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar serviços"})
			return
		}

		if err := db.Where("supplier_link_id = ?", existingSupplier.ID).Delete(&schemas.SupplierProduct{}).Error; err != nil {
			log.Printf("Erro ao remover produtos existentes: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar produtos"})
			return
		}

		// Adiciona os novos serviços e produtos
		if err := createServicesAndProducts(existingSupplier.ID, input.ServiceIDs, input.ProductIDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Fornecedor reativado com sucesso", "supplier": existingSupplier})
		return
	}

	// Se não existe, busca no banco externo e cria novo
	fornecedor, err := getFornecedorByCNPJ(input.SupplierCNPJ)
	if err != nil {
		log.Printf("Erro ao buscar fornecedor no banco externo: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Fornecedor não encontrado no banco externo"})
		return
	}

	// Cria novo fornecedor
	supplierLink := schemas.SupplierLink{
		CNPJ:       input.SupplierCNPJ,
		CategoryID: input.CategoryID,
	}

	if err := db.Create(&supplierLink).Error; err != nil {
		log.Printf("Erro ao criar fornecedor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar fornecedor"})
		return
	}

	// Cria os serviços e produtos associados
	if err := createServicesAndProducts(supplierLink.ID, input.ServiceIDs, input.ProductIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Fornecedor cadastrado com sucesso",
		"supplier": fornecedor,
	})
}

// Função auxiliar para criar serviços e produtos
func createServicesAndProducts(supplierID uint, serviceIDs []uint, productIDs []uint) error {
	// Verifica se os serviços existem
	for _, serviceID := range serviceIDs {
		var service schemas.Service
		if err := db.First(&service, serviceID).Error; err != nil {
			log.Printf("Serviço ID %d não encontrado: %v", serviceID, err)
			return fmt.Errorf("serviço ID %d não encontrado", serviceID)
		}
		
		supplierService := schemas.SupplierService{
			SupplierLinkID: supplierID,
			ServiceID:      serviceID,
		}
		if err := db.Create(&supplierService).Error; err != nil {
			log.Printf("Erro ao criar serviço para fornecedor: %v", err)
			return fmt.Errorf("erro ao vincular serviço")
		}
	}

	// Verifica se os produtos existem
	for _, productID := range productIDs {
		var product schemas.Product
		if err := db.First(&product, productID).Error; err != nil {
			log.Printf("Produto ID %d não encontrado: %v", productID, err)
			return fmt.Errorf("produto ID %d não encontrado", productID)
		}

		supplierProduct := schemas.SupplierProduct{
			SupplierLinkID: supplierID,
			ProductID:      productID,
		}
		if err := db.Create(&supplierProduct).Error; err != nil {
			log.Printf("Erro ao criar produto para fornecedor: %v", err)
			return fmt.Errorf("erro ao vincular produto")
		}
	}

	return nil
}

// getFornecedorByCNPJ busca um fornecedor externo pelo CNPJ
func getFornecedorByCNPJ(cnpj string) (*schemas.ExternalSupplier, error) {
	fornecedores, err := getFornecedoresExternosFromDatabase()
	if err != nil {
		log.Printf("Erro ao buscar fornecedores do banco de dados: %v", err)
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

	log.Printf("Fornecedor não encontrado para o CNPJ: %s", cnpj)
	return nil, fmt.Errorf("Fornecedor não encontrado")
}

