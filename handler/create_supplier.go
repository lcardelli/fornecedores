package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

// Cria um novo fornecedor
func CreateSupplierHandler(c *gin.Context) {
	log.Println("Iniciando CreateSupplierHandler")
	// Estrutura para receber os dados do fornecedor
	var input struct {
		CNPJ       string   `json:"supplier_cnpj" binding:"required"`
		CategoryID string   `json:"category_id" binding:"required"`
		ServiceIDs []string `json:"service_ids" binding:"required,min=1"`
	}

	// Faz o bind do JSON recebido para a estrutura
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Erro ao fazer bind dos dados de entrada: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	// Log dos dados recebidos
	log.Printf("Dados recebidos: CNPJ=%s, CategoryID=%s, ServiceIDs=%v", input.CNPJ, input.CategoryID, input.ServiceIDs)

	// Converter CategoryID para uint
	categoryID, err := strconv.ParseUint(input.CategoryID, 10, 32)
	if err != nil {
		log.Printf("Erro ao converter CategoryID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryID inválido"})
		return
	}

	// Buscar os IDs dos serviços pelo nome
	var serviceIDs []uint
	for _, serviceName := range input.ServiceIDs {
		var service schemas.Service
		if err := db.Where("name = ?", serviceName).First(&service).Error; err != nil {
			log.Printf("Erro ao buscar serviço pelo nome: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Serviço não encontrado: " + serviceName})
			return
		}
		serviceIDs = append(serviceIDs, service.ID)
	}

	// Verificar se o fornecedor já está vinculado
	var existingLink schemas.SupplierLink
	if err := db.Where("cnpj = ?", input.CNPJ).First(&existingLink).Error; err == nil {
		log.Printf("Fornecedor já vinculado: CNPJ=%s", input.CNPJ)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fornecedor já vinculado"})
		return
	}

	// Buscar o fornecedor no banco de dados externo
	fornecedor, err := getFornecedorByCNPJ(input.CNPJ)
	if err != nil {
		log.Printf("Erro ao buscar fornecedor no banco de dados externo: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Fornecedor não encontrado no banco de dados externo"})
		return
	}

	log.Printf("Fornecedor encontrado no banco de dados externo: %+v", fornecedor)

	// Criar o vínculo do fornecedor
	supplierLink := schemas.SupplierLink{
		CNPJ:       input.CNPJ,
		CategoryID: uint(categoryID),
	}

	// Salvar o vínculo do fornecedor no banco de dados
	if err := db.Create(&supplierLink).Error; err != nil {
		log.Printf("Erro ao criar vínculo do fornecedor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar vínculo do fornecedor"})
		return
	}

	log.Printf("Vínculo do fornecedor criado com sucesso: ID=%d", supplierLink.ID)

	// Vincular serviços
	for _, serviceID := range serviceIDs {
		supplierService := schemas.SupplierService{
			SupplierLinkID: supplierLink.ID,
			ServiceID:      serviceID,
		}
		if err := db.Create(&supplierService).Error; err != nil {
			log.Printf("Erro ao vincular serviço: SupplierLinkID=%d, ServiceID=%d, Erro=%v", supplierLink.ID, serviceID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao vincular serviço"})
			return
		}
		log.Printf("Serviço vinculado com sucesso: SupplierLinkID=%d, ServiceID=%d", supplierLink.ID, serviceID)
	}

	log.Printf("Fornecedor vinculado com sucesso: CNPJ=%s", input.CNPJ)
	c.JSON(http.StatusCreated, gin.H{"message": "Fornecedor vinculado com sucesso", "supplier": fornecedor})
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
