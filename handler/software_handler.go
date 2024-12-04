package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

// SoftwareInput estrutura para receber os dados do frontend
type SoftwareInput struct {
	Name        string `json:"name"`
	Publisher   string `json:"publisher"`
	Description string `json:"description"`
}

// RenderManageSoftwareHandler renderiza a página de gerenciamento de softwares
func RenderManageSoftwareHandler(c *gin.Context) {
	var softwares []schemas.Software
	var publishers []string

	// Carrega os softwares com suas licenças
	if err := db.Preload("Licenses").Find(&softwares).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar softwares",
		})
		return
	}

	// Obtém lista única de fabricantes
	if err := db.Model(&schemas.Software{}).
		Distinct("publisher").
		Pluck("publisher", &publishers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar fabricantes",
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

	c.HTML(http.StatusOK, "manage_software.html", gin.H{
		"softwares": softwares,
		"publishers": publishers,
		"user":      currentUser,
	})
}

// CreateSoftwareHandler cria um novo software
func CreateSoftwareHandler(c *gin.Context) {
	var input SoftwareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	software := schemas.Software{
		Name:        input.Name,
		Publisher:   input.Publisher,
		Description: input.Description,
	}

	if err := db.Create(&software).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar software"})
		return
	}

	c.JSON(http.StatusCreated, software)
}

// UpdateSoftwareHandler atualiza um software existente
func UpdateSoftwareHandler(c *gin.Context) {
	id := c.Param("id")

	// Converte o ID para uint
	softwareID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Busca o software existente
	var software schemas.Software
	if err := db.First(&software, softwareID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software não encontrado"})
		return
	}

	// Recebe os dados da atualização
	var input SoftwareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza apenas os campos necessários
	software.Name = input.Name
	software.Publisher = input.Publisher
	software.Description = input.Description

	// Salva as alterações
	if err := db.Save(&software).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar software"})
		return
	}

	c.JSON(http.StatusOK, software)
}

// DeleteSoftwareHandler deleta um software
func DeleteSoftwareHandler(c *gin.Context) {
	id := c.Param("id")

	// Verifica se existem licenças associadas
	var count int64
	if err := db.Model(&schemas.License{}).Where("software_id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar licenças"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível excluir software com licenças associadas"})
		return
	}

	if err := db.Delete(&schemas.Software{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar software"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Software deletado com sucesso"})
}

// GetSoftwareHandler retorna um software específico
func GetSoftwareHandler(c *gin.Context) {
	id := c.Param("id")
	var software schemas.Software

	if err := db.First(&software, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software não encontrado"})
		return
	}

	c.JSON(http.StatusOK, software)
}
