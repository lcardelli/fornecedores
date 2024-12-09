package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"github.com/lcardelli/fornecedores/utils"
	"gorm.io/gorm"
)

// Adicione esta função helper

// RenderManageLicensesHandler renderiza a página de gerenciamento de licenças
func RenderManageLicensesHandler(c *gin.Context) {
	var licenses []schemas.License
	var softwares []schemas.Software
	var users []schemas.User
	var periodRenews []schemas.PeriodRenew
	var departments []schemas.Departament
	var totalCost float64
	var years []string

	// Buscar anos únicos das datas de expiração
	if err := db.Table("licenses").
		Select("DISTINCT YEAR(expiry_date) as year").
		Where("expiry_date IS NOT NULL").
		Where("licenses.deleted_at IS NULL").
		Where("expiry_date > ?", "1000-01-01").
		Order("year DESC").
		Pluck("year", &years).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar anos das licenças",
		})
		return
	}

	// Carrega as licenças com seus relacionamentos
	if err := db.Preload("Software").
		Preload("Status").
		Preload("AssignedUsers").
		Preload("PeriodRenew").
		Preload("Department").
		Find(&licenses).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar licenças",
		})
		return
	}

	// Carrega os períodos de renovação
	if err := db.Find(&periodRenews).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar períodos de renovação",
		})
		return
	}

	// Calcula o total
	for _, license := range licenses {
		if license.Cost > 0 {
			totalCost += license.Cost
		}
	}

	// Carrega a lista de softwares para o select
	if err := db.Model(&schemas.Software{}).
		Where("deleted_at IS NULL").
		Order("name ASC").
		Find(&softwares).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar softwares",
		})
		return
	}

	// Carrega a lista de usuários para o select
	if err := db.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar usuários",
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

	// Antes de renderizar o template, formate o custo total
	formattedTotalCost := utils.FormatMoney(totalCost)

	// Carrega os departamentos
	if err := db.Find(&departments).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar departamentos",
		})
		return
	}

	c.HTML(http.StatusOK, "manage_licenses.html", gin.H{
		"licenses":     licenses,
		"softwares":    softwares,
		"users":        users,
		"periodRenews": periodRenews,
		"departments":  departments,
		"user":         currentUser,
		"totalCost":    formattedTotalCost,
		"formatMoney":  utils.FormatMoney, // Adiciona a função helper ao template
		"years":        years,
	})
}

// CreateLicenseHandler cria uma nova licença
func CreateLicenseHandler(c *gin.Context) {
	var input schemas.License
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Busca o status "Ativa" do banco
	var status schemas.Status
	if err := db.Where("name = ?", "Ativa").First(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar status"})
		return
	}

	input.StatusID = status.ID

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar licença"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// DeleteLicenseHandler deleta uma licença
func DeleteLicenseHandler(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&schemas.License{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar licença"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licença deletada com sucesso"})
}

// GetLicense busca uma licença específica
func GetLicense(c *gin.Context) {
	id := c.Param("id")

	var license schemas.License
	result := db.Preload("Software").
		Preload("Status").
		Preload("AssignedUsers").
		First(&license, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Licença não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar licença"})
		return
	}

	// Recalcula o status
	if err := license.CalculateStatus(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao calcular status"})
		return
	}

	// Atualiza o status no banco
	if err := db.Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
		return
	}

	c.JSON(http.StatusOK, license)
}

// UpdateLicenseHandler atualiza uma licença existente
func UpdateLicenseHandler(c *gin.Context) {
	id := c.Param("id")

	var license schemas.License
	if err := db.First(&license, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Licença não encontrada"})
		return
	}

	var input schemas.License
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Se apenas o campo blocked foi enviado, atualiza somente ele
	if input.Blocked != license.Blocked && input.Software.ID == 0 && input.LicenseKey == "" {
		if err := db.Model(&license).Update("blocked", input.Blocked).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status da licença"})
			return
		}
		// Recarrega a licença com os relacionamentos
		if err := db.Preload("Software").Preload("Status").Preload("AssignedUsers").First(&license, license.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recarregar licença"})
			return
		}
		c.JSON(http.StatusOK, license)
		return
	}

	// Validar department_id
	if input.DepartmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Departamento é obrigatório"})
		return
	}

	// Verificar se o departamento existe
	var department schemas.Departament
	if err := db.First(&department, input.DepartmentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Departamento inválido"})
		return
	}

	// Se period_renew_id for 0, define como nil
	if input.PeriodRenewID != nil && *input.PeriodRenewID == 0 {
		input.PeriodRenewID = nil
	}

	// Atualiza os campos da licença
	if err := db.Model(&license).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar licença"})
		return
	}

	// Atualiza os usuários designados
	if err := db.Model(&license).Association("AssignedUsers").Replace(input.AssignedUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuários designados"})
		return
	}

	// Recalcula o status após a atualização
	if err := license.CalculateStatus(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao calcular status"})
		return
	}

	// Salva o novo status
	if err := db.Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
		return
	}

	// Recarrega a licença com todos os relacionamentos
	if err := db.Preload("Software").
		Preload("Status").
		Preload("AssignedUsers").
		First(&license, license.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recarregar licença"})
		return
	}

	c.JSON(http.StatusOK, license)
}

// RenderViewLicensesPage renderiza a página de visualização de licenças
func RenderViewLicensesPage(c *gin.Context) {
	var departments []schemas.Departament
	var years []string
	var softwares []schemas.Software

	// Buscar anos únicos das datas de expiração
	if err := db.Table("licenses").
		Select("DISTINCT YEAR(expiry_date) as year").
		Where("expiry_date IS NOT NULL").
		Where("licenses.deleted_at IS NULL").
		Where("expiry_date > ?", "1000-01-01").
		Order("year DESC").
		Pluck("year", &years).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar anos das licenças",
		})
		return
	}

	// Carregar departamentos
	if err := db.Find(&departments).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar departamentos",
		})
		return
	}

	// Carregar softwares
	if err := db.Table("licenses").
		Joins("JOIN softwares ON licenses.software_id = softwares.id").
		Select("DISTINCT softwares.id, softwares.name").
		Where("licenses.deleted_at IS NULL").
		Order("softwares.name ASC").
		Scan(&softwares).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar softwares",
		})
		return
	}

	RenderTemplate(c, "list_licenses.html", gin.H{
		"activeMenu":  "visualizar-licencas",
		"years":       years,
		"departments": departments,
		"softwares":   softwares,
	})
}

// ListLicensesHandler processa a requisição de listagem de licenças
func ListLicensesHandler(c *gin.Context) {
	// Obter parâmetros de filtro
	softwareID := c.Query("software_id")
	statusID := c.Query("status_id")
	dateFilter := c.Query("date")
	departmentFilter := c.Query("department")
	typeFilter := c.Query("type")

	// Primeiro, buscar todas as licenças para atualizar seus status
	var allLicenses []schemas.License
	query := db.Preload("Software").
		Preload("Status").
		Preload("PeriodRenew").
		Preload("Department").
		Where("licenses.deleted_at IS NULL")

	if err := query.Find(&allLicenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar licenças"})
		return
	}

	// Atualizar o status de todas as licenças
	for i := range allLicenses {
		if err := allLicenses[i].CalculateStatus(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao calcular status das licenças"})
			return
		}
		// Salvar o status atualizado no banco
		if err := db.Save(&allLicenses[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar status atualizado"})
			return
		}
	}

	// Agora criar uma nova query com os filtros
	filteredQuery := db.Preload("Software").
		Preload("Status").
		Preload("PeriodRenew").
		Preload("Department").
		Where("licenses.deleted_at IS NULL")

	// Filtro por software usando ID
	if softwareID != "" {
		softwareIDInt, err := strconv.Atoi(softwareID)
		if err == nil {
			filteredQuery = filteredQuery.Where("licenses.software_id = ?", softwareIDInt)
		}
	}

	// Filtro por tipo
	if typeFilter != "" {
		filteredQuery = filteredQuery.Where("licenses.type = ?", typeFilter)
	}

	// Filtro por status usando ID como número
	if statusID != "" {
		statusIDInt, err := strconv.Atoi(statusID)
		if err == nil {
			filteredQuery = filteredQuery.Where("licenses.status_id = ?", statusIDInt)
		}
	}

	// Filtro por ano de expiração
	if dateFilter != "" {
		filteredQuery = filteredQuery.Where("YEAR(licenses.expiry_date) = ?", dateFilter)
	}

	// Filtro por departamento
	if departmentFilter != "" {
		filteredQuery = filteredQuery.Where("licenses.department_id = ?", departmentFilter)
	}

	// Buscar licenças filtradas
	var filteredLicenses []schemas.License
	if err := filteredQuery.Find(&filteredLicenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar licenças filtradas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"licenses": filteredLicenses,
	})
}

// GetFilteredLicenses retorna as licenças filtradas de acordo com os parâmetros
func GetFilteredLicenses(search, status, dateFilter string) []schemas.License {
	query := db.Table("licenses").
		Preload("Software").
		Preload("Status")

	// Aplicar filtros
	if search != "" {
		query = query.Where(
			"license_key LIKE ? OR softwares.name LIKE ?",
			"%"+search+"%", "%"+search+"%",
		).Joins("LEFT JOIN softwares ON licenses.software_id = softwares.id")
	}

	// Filtro por status usando ID
	if status != "" {
		query = query.Where("licenses.status_id = ?", status)
	}

	// Filtro por ano de expiração
	if dateFilter != "" {
		query = query.Where("YEAR(licenses.expiry_date) = ?", dateFilter)
	}

	var licenses []schemas.License
	query.Find(&licenses)

	return licenses
}

// Adicione esta estrutura para receber os IDs
type BatchDeleteRequest struct {
	IDs []uint `json:"ids"`
}

// Adicione este handler para deleção em lote
func DeleteBatchLicenses(c *gin.Context) {
	var request BatchDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs inválidos"})
		return
	}

	if len(request.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum ID fornecido"})
		return
	}

	// Usar o DB para deletar múltiplas licenças
	result := db.Where("id IN (?)", request.IDs).Delete(&schemas.License{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir licenças"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licenças excluídas com sucesso"})
}
