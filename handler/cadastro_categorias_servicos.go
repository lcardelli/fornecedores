package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
	"gorm.io/gorm"
)

func CadastroCategoriaHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do usuário"})
		return
	}

	c.HTML(http.StatusOK, "cadastro_categoria.html", gin.H{
		"user":       user,
		"activeMenu": "cadastro-categoria",
	})
}

func CreateCategoryHandler(c *gin.Context) {
	var category schemas.SupplierCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar categoria"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	var category schemas.SupplierCategory

	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoria não encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategoryHandler(c *gin.Context) {
	id := c.Param("id")

	// Verifica se o ID é válido
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria não fornecido"})
		return
	}

	var category schemas.SupplierCategory
	result := db.First(&category, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Categoria não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categoria"})
		}
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar categoria"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria deletada com sucesso"})
}

func ListCategoriesHandler(c *gin.Context) {
	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categorias"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func CadastroServicoHandler(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do usuário"})
		return
	}

	var categories []schemas.SupplierCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categorias"})
		return
	}

	c.HTML(http.StatusOK, "cadastro_servico.html", gin.H{
		"user":       user,
		"Categories": categories,
		"activeMenu": "cadastro-servico",
	})
}

func CreateServiceHandler(c *gin.Context) {
	log.Println("Iniciando CreateServiceHandler")

	var service schemas.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Dados do serviço recebidos: %+v", service)

	if err := db.Create(&service).Error; err != nil {
		log.Printf("Erro ao criar serviço no banco de dados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar serviço"})
		return
	}

	log.Printf("Serviço criado com sucesso: ID %d", service.ID)
	c.JSON(http.StatusCreated, service)
}
