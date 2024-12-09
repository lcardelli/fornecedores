package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"github.com/lcardelli/fornecedores/utils"
)

// Adicione esta constante
const uploadDir = "uploads/contracts"

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

	formattedTotalValue := utils.FormatMoney(totalValue)

	// Buscar anos únicos das datas dos contratos
	var years []string
	if err := db.Raw(`
		SELECT DISTINCT YEAR(initial_date) as year FROM contracts 
		WHERE initial_date IS NOT NULL AND deleted_at IS NULL
		UNION 
		SELECT DISTINCT YEAR(final_date) as year FROM contracts 
		WHERE final_date IS NOT NULL AND deleted_at IS NULL
		ORDER BY year ASC
	`).Pluck("year", &years).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar anos dos contratos",
		})
		return
	}

	c.HTML(http.StatusOK, "manage_contracts.html", gin.H{
		"contracts":             contracts,
		"departments":           departments,
		"branches":              branches,
		"costCenters":           costCenters,
		"contractStatuses":      contractStatuses,
		"terminationConditions": terminationConditions,
		"user":                  currentUser,
		"totalValue":            formattedTotalValue,
		"formatMoney":           utils.FormatMoney,
		"years":                 years,
		"activeMenu":            "contratos",
	})
}

func CreateContractHandler(c *gin.Context) {
	// Garante que o diretório de upload existe
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar diretório de uploads"})
		return
	}

	// Parse do formulário multipart
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário: " + err.Error()})
		return
	}

	var input schemas.Contract

	// Bind dos campos do formulário
	input.Name = c.PostForm("name")
	input.ContractNumber = c.PostForm("contract_number")
	input.Object = c.PostForm("object")

	// Converter valor
	if value, err := strconv.ParseFloat(c.PostForm("value"), 64); err == nil {
		input.Value = value
	}

	// Converter IDs
	if departmentID, err := strconv.ParseUint(c.PostForm("department_id"), 10, 32); err == nil {
		input.DepartmentID = uint(departmentID)
	}
	if branchID, err := strconv.ParseUint(c.PostForm("branch_id"), 10, 32); err == nil {
		input.BranchID = uint(branchID)
	}
	if costCenterID, err := strconv.ParseUint(c.PostForm("cost_center_id"), 10, 32); err == nil {
		input.CostCenterID = uint(costCenterID)
	}
	if terminationConditionID, err := strconv.ParseUint(c.PostForm("termination_condition_id"), 10, 32); err == nil {
		input.TerminationConditionID = uint(terminationConditionID)
	}
	if statusID, err := strconv.ParseUint(c.PostForm("status_id"), 10, 32); err == nil {
		input.StatusID = uint(statusID)
	}

	// Converter datas
	if initialDate, err := time.Parse("2006-01-02T15:04:05Z", c.PostForm("initial_date")); err == nil {
		input.InitialDate = initialDate
	}
	if finalDate, err := time.Parse("2006-01-02T15:04:05Z", c.PostForm("final_date")); err == nil {
		input.FinalDate = finalDate
	}

	input.Observations = c.PostForm("notes")

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

	// Iniciar transação
	tx := db.Begin()

	// Criar contrato
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar contrato"})
		return
	}

	// Processar arquivos
	form, _ := c.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		// Criar nome único para o arquivo
		filename := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", input.ID, file.Filename))

		// Salvar arquivo no sistema
		if err := c.SaveUploadedFile(file, filename); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
			return
		}

		// Criar registro do anexo
		anexo := schemas.ContractAnexo{
			ContractID: input.ID,
			Name:       file.Filename,
			Path:       filename,
			FileType:   file.Header.Get("Content-Type"),
			FileSize:   file.Size,
		}

		if err := tx.Create(&anexo).Error; err != nil {
			tx.Rollback()
			// Remover arquivo salvo em caso de erro
			os.Remove(filename)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar anexo"})
			return
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar transação"})
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

	// Parse do formulário multipart
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário: " + err.Error()})
		return
	}

	// Criar uma nova estrutura para as atualizações
	var updates = make(map[string]interface{})

	// Atualizar campos básicos
	if name := c.PostForm("name"); name != "" {
		updates["name"] = name
	}
	if contractNumber := c.PostForm("contract_number"); contractNumber != "" {
		updates["contract_number"] = contractNumber
	}
	if object := c.PostForm("object"); object != "" {
		updates["object"] = object
	}
	if notes := c.PostForm("notes"); notes != "" {
		updates["observations"] = notes
	}

	// Converter e atualizar valor
	if valueStr := c.PostForm("value"); valueStr != "" {
		if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
			updates["value"] = value
		}
	}

	// Converter e atualizar IDs
	if departmentID := c.PostForm("department_id"); departmentID != "" {
		if id, err := strconv.ParseUint(departmentID, 10, 32); err == nil {
			updates["department_id"] = uint(id)
		}
	}
	if branchID := c.PostForm("branch_id"); branchID != "" {
		if id, err := strconv.ParseUint(branchID, 10, 32); err == nil {
			updates["branch_id"] = uint(id)
		}
	}
	if costCenterID := c.PostForm("cost_center_id"); costCenterID != "" {
		if id, err := strconv.ParseUint(costCenterID, 10, 32); err == nil {
			updates["cost_center_id"] = uint(id)
		}
	}
	if terminationConditionID := c.PostForm("termination_condition_id"); terminationConditionID != "" {
		if id, err := strconv.ParseUint(terminationConditionID, 10, 32); err == nil {
			updates["termination_condition_id"] = uint(id)
		}
	}
	if statusID := c.PostForm("status_id"); statusID != "" {
		if id, err := strconv.ParseUint(statusID, 10, 32); err == nil {
			updates["status_id"] = uint(id)
		}
	}

	// Converter e atualizar datas
	if initialDate := c.PostForm("initial_date"); initialDate != "" {
		if date, err := time.Parse("2006-01-02T15:04:05Z", initialDate); err == nil {
			updates["initial_date"] = date
		}
	}
	if finalDate := c.PostForm("final_date"); finalDate != "" {
		if date, err := time.Parse("2006-01-02T15:04:05Z", finalDate); err == nil {
			updates["final_date"] = date
		}
	}

	// Obter usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	currentUser := userInterface.(schemas.User)

	// Atualizar campos de auditoria
	updates["updated_by"] = currentUser.ID
	updates["last_modified"] = time.Now()

	// Iniciar transação
	tx := db.Begin()

	// Atualizar contrato
	if err := tx.Model(&existingContract).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar contrato"})
		return
	}

	// Processar novos arquivos
	form, _ := c.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		// Criar nome único para o arquivo
		filename := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", existingContract.ID, file.Filename))

		// Salvar arquivo no sistema
		if err := c.SaveUploadedFile(file, filename); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
			return
		}

		// Criar registro do anexo
		anexo := schemas.ContractAnexo{
			ContractID: existingContract.ID,
			Name:       file.Filename,
			Path:       filename,
			FileType:   file.Header.Get("Content-Type"),
			FileSize:   file.Size,
		}

		if err := tx.Create(&anexo).Error; err != nil {
			tx.Rollback()
			// Remover arquivo salvo em caso de erro
			os.Remove(filename)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar anexo"})
			return
		}
	}

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar transação"})
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

	// Primeiro verifica se está vencido
	if now.After(contract.FinalDate) {
		contract.StatusID = 3 // Vencido
		return nil
	}

	// Calcula a diferença em dias
	daysUntilExpiration := contract.FinalDate.Sub(now).Hours() / 24

	// Se faltam 30 dias ou menos, está próximo ao vencimento
	if daysUntilExpiration <= 30 && daysUntilExpiration >= 0 {
		contract.StatusID = 2 // Próximo do Vencimento
		return nil
	}

	contract.StatusID = 1 // Em Vigor
	return nil
}

func DownloadContractAttachmentHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var anexo schemas.ContractAnexo
	if err := db.First(&anexo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anexo não encontrado"})
		return
	}

	// Verifica se o arquivo existe
	if _, err := os.Stat(anexo.Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	// Define o nome do arquivo para download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", anexo.Name))
	c.Header("Content-Type", anexo.FileType)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve o arquivo
	c.File(anexo.Path)
}

// RenderListContractsHandler renderiza a página de visualização de contratos
func RenderListContractsHandler(c *gin.Context) {
	var contracts []schemas.Contract
	var departments []schemas.ContractDepartament
	var branches []schemas.ContractFilial
	var costCenters []schemas.ContractCentroCusto
	var contractStatuses []schemas.ContractStatus
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

	// Buscar anos únicos das datas dos contratos
	var years []string
	if err := db.Raw(`
		SELECT DISTINCT YEAR(initial_date) as year FROM contracts 
		WHERE initial_date IS NOT NULL AND deleted_at IS NULL
		UNION 
		SELECT DISTINCT YEAR(final_date) as year FROM contracts 
		WHERE final_date IS NOT NULL AND deleted_at IS NULL
		ORDER BY year ASC
	`).Pluck("year", &years).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar anos dos contratos",
		})
		return
	}

	formattedTotalValue := utils.FormatMoney(totalValue)

	c.HTML(http.StatusOK, "list_contracts.html", gin.H{
		"contracts":             contracts,
		"departments":           departments,
		"branches":              branches,
		"costCenters":          costCenters,
		"contractStatuses":      contractStatuses,
		"user":                 currentUser,
		"totalValue":           formattedTotalValue,
		"formatMoney":          utils.FormatMoney,
		"years":                years,
		"activeMenu":           "contratos",
	})
}
