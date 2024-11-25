package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/gorm"
)

// RenderManageLicensesHandler renderiza a página de gerenciamento de licenças
func RenderManageLicensesHandler(c *gin.Context) {
	var licenses []schemas.License
	var softwares []schemas.Software
	var users []schemas.User
	var totalCost float64

	// Carrega as licenças com seus relacionamentos
	if err := db.Preload("Software").Preload("AssignedUsers").Find(&licenses).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Erro ao carregar licenças",
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
	if err := db.Find(&softwares).Error; err != nil {
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

	c.HTML(http.StatusOK, "manage_licenses.html", gin.H{
		"licenses":  licenses,
		"softwares": softwares,
		"users":     users,
		"user":      currentUser,
		"totalCost": totalCost,
	})
}

// CreateLicenseHandler cria uma nova licença
func CreateLicenseHandler(c *gin.Context) {
	var input schemas.License
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define o status inicial como "Ativa"
	input.Status = "Ativa"

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

	c.JSON(http.StatusOK, license)
}
