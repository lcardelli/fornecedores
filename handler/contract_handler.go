package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// RenderManageContractsHandler renderiza a página de gerenciamento de contratos
func RenderManageContractsHandler(c *gin.Context) {
	var contracts []schemas.Contract
	var departments []schemas.ContractDepartament
	var branches []schemas.ContractFilial
	var costCenters []schemas.ContractCentroCusto
	var contractStatuses []schemas.ContractStatus
	var terminationConditions []schemas.ContractCondicaoRescisao
	var totalValue float64

	// Carrega os contratos com seus relacionamentos
	if err := db.Preload("Status").
		Preload("CostCenter").
		Preload("Branch").
		Preload("Department").
		Preload("TerminationCondition").
		Preload("Attachments").
		Find(&contracts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar contratos",
		})
		return
	}

	// Calcula o valor total
	for _, contract := range contracts {
		if contract.Value > 0 {
			totalValue += contract.Value
		}
	}

	// Carrega departamentos
	if err := db.Find(&departments).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar departamentos",
		})
		return
	}

	// Carrega filiais
	if err := db.Find(&branches).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar filiais",
		})
		return
	}

	// Carrega centros de custo
	if err := db.Find(&costCenters).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar centros de custo",
		})
		return
	}

	// Carrega status
	if err := db.Find(&contractStatuses).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar status",
		})
		return
	}

	// Carrega condições de rescisão
	if err := db.Find(&terminationConditions).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar condições de rescisão",
		})
		return
	}

	// Obtém o usuário atual do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Usuário não encontrado no contexto",
		})
		return
	}
	currentUser, ok := userInterface.(schemas.User)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao processar dados do usuário",
		})
		return
	}

	// Atualiza o status de cada contrato
	for i := range contracts {
		if err := updateContractStatus(&contracts[i]); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "Erro ao atualizar status dos contratos",
			})
			return
		}
		
		// Salva o novo status
		if err := db.Save(&contracts[i]).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "Erro ao salvar status dos contratos",
			})
			return
		}
	}

	formattedTotalValue := formatMoney(totalValue)

	c.HTML(http.StatusOK, "manage_contracts.html", gin.H{
		"contracts":   contracts,
		"departments": departments,
		"branches":    branches,
		"costCenters": costCenters,
		"contractStatuses": contractStatuses,
		"terminationConditions": terminationConditions,
		"user":        currentUser,
		"totalValue":  formattedTotalValue,
		"formatMoney": formatMoney,
		"activeMenu":  "contratos",
	})
}

func CreateContractHandler(c *gin.Context) {
	var input schemas.Contract
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obter usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	currentUser := userInterface.(schemas.User)

	// Adicionar campos de auditoria
	input.CreatedBy = currentUser.ID
	input.UpdatedBy = currentUser.ID
	input.LastModified = time.Now()

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar contrato"})
		return
	}

	// Recarregar o contrato com todos os relacionamentos
	if err := db.Preload("Status").
		Preload("CostCenter").
		Preload("Branch").
		Preload("Department").
		Preload("TerminationCondition").
		Preload("Attachments").
		First(&input, input.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recarregar contrato"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func UpdateContractHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o contrato existe
	var existingContract schemas.Contract
	if err := db.First(&existingContract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contrato não encontrado"})
		return
	}

	var input schemas.Contract
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obter usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	currentUser := userInterface.(schemas.User)

	// Atualizar campos de auditoria
	input.UpdatedBy = currentUser.ID
	input.LastModified = time.Now()

	if err := db.Model(&existingContract).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar contrato"})
		return
	}

	// Recarregar o contrato com todos os relacionamentos
	if err := db.Preload("Status").
		Preload("CostCenter").
		Preload("Branch").
		Preload("Department").
		Preload("TerminationCondition").
		Preload("Attachments").
		First(&existingContract, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recarregar contrato"})
		return
	}

	c.JSON(http.StatusOK, existingContract)
}

func DeleteContractHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := db.Delete(&schemas.Contract{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar contrato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contrato deletado com sucesso"})
}

// Estrutura para deleção em lote
type BatchDeleteContractsRequest struct {
	IDs []uint `json:"ids"`
}

func DeleteBatchContracts(c *gin.Context) {
	var request BatchDeleteContractsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	if len(request.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum ID fornecido"})
		return
	}

	if err := db.Where("id IN (?)", request.IDs).Delete(&schemas.Contract{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir contratos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contratos excluídos com sucesso"})
}

func GetContractHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var contract schemas.Contract
	if err := db.Preload("Status").
		Preload("CostCenter").
		Preload("Branch").
		Preload("Department").
		Preload("TerminationCondition").
		Preload("Attachments").
		First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contrato não encontrado"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func GetAllContractsHandler(c *gin.Context) {
	var contracts []schemas.Contract
	if err := db.Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar contratos"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

func GetContractAditivosHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var aditivos []schemas.ContractAditivo
	if err := db.Where("contract_id = ?", id).Find(&aditivos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar aditivos"})
		return
	}

	c.JSON(http.StatusOK, aditivos)
}

func CreateContractAditivoHandler(c *gin.Context) {
	var input schemas.ContractAditivo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar aditivo"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func UpdateContractAditivoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input schemas.ContractAditivo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&input).Where("id = ?", id).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar aditivo"})
		return
	}

	c.JSON(http.StatusOK, input)
}

func DeleteContractAditivoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := db.Delete(&schemas.ContractAditivo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar aditivo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Aditivo deletado com sucesso"})
}

func GetAllContractAditivosHandler(c *gin.Context) {
	var aditivos []schemas.ContractAditivo
	if err := db.Find(&aditivos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar aditivos"})
		return
	}

	c.JSON(http.StatusOK, aditivos)
}

// Adicionar esta função auxiliar
func updateContractStatus(contract *schemas.Contract) error {
	// Verifica se tem aditivos
	var aditivosCount int64
	if err := db.Model(&schemas.ContractAditivo{}).Where("contract_id = ?", contract.ID).Count(&aditivosCount).Error; err != nil {
		return err
	}

	if aditivosCount > 0 {
		contract.StatusID = 4 // Renovado por Aditivo
		return nil
	}

	now := time.Now()
	if now.After(contract.FinalDate) {
		contract.StatusID = 3 // Vencido
		return nil
	}

	daysUntilExpiration := contract.FinalDate.Sub(now).Hours() / 24
	if daysUntilExpiration <= 30 {
		contract.StatusID = 2 // Próximo ao Vencimento
		return nil
	}

	contract.StatusID = 1 // Em Vigor
	return nil
}
