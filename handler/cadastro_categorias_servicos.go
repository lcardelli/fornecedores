package handler

import (
	"log"
	"net/http"
	"strconv"

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

	var services []schemas.Service
	if err := db.Preload("Category").Find(&services).Error; err != nil {
		log.Printf("Erro ao buscar serviços: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	c.HTML(http.StatusOK, "cadastro_servico.html", gin.H{
		"user":       user,
		"Categories": categories,
		"Services":   services,
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

func ListServicesHandler(c *gin.Context) {
	var services []schemas.Service
	if err := db.Preload("Category").Find(&services).Error; err != nil {
		log.Printf("Erro ao buscar serviços: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
		return
	}

	// Criar uma slice de ServiceResponse
	var serviceResponses []schemas.ServiceResponse
	for _, service := range services {
		serviceResponse := schemas.ServiceResponse{
			ID:         service.ID,
			Name:       service.Name,
			CategoryID: service.CategoryID,
			Category: schemas.SupplierCategoryResponse{
				ID:   service.Category.ID,
				Name: service.Category.Name,
			},
		}
		serviceResponses = append(serviceResponses, serviceResponse)
	}

	c.JSON(http.StatusOK, serviceResponses)
}

func UpdateServiceHandler(c *gin.Context) {
	id := c.Param("id")
	log.Printf("ID recebido para atualização: %s", id)

	var service schemas.Service
	if err := db.Preload("Category").First(&service, id).Error; err != nil {
		log.Printf("Erro ao buscar serviço: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	var updateData struct {
		Name       string `json:"name"`
		CategoryID uint   `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Name = updateData.Name
	service.CategoryID = updateData.CategoryID

	if err := db.Save(&service).Error; err != nil {
		log.Printf("Erro ao atualizar serviço: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar serviço"})
		return
	}

	// Recarrega o serviço com a categoria atualizada
	if err := db.Preload("Category").First(&service, id).Error; err != nil {
		log.Printf("Erro ao recarregar serviço atualizado: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recarregar serviço atualizado"})
		return
	}

	log.Printf("Serviço atualizado com sucesso: %+v", service)
	c.JSON(http.StatusOK, service)
}

func DeleteServiceHandler(c *gin.Context) {
	id := c.Param("id")
	log.Printf("ID recebido para deleção: %s", id)

	if id == "" || id == "undefined" {
		log.Println("ID do serviço não fornecido ou inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço não fornecido ou inválido"})
		return
	}

	// Converte o ID para uint
	serviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("Erro ao converter ID para uint: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	log.Printf("ID convertido para uint: %d", serviceID)

	var service schemas.Service
	result := db.First(&service, uint(serviceID))

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Serviço não encontrado para o ID: %d", serviceID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		} else {
			log.Printf("Erro ao buscar serviço: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviço"})
		}
		return
	}

	log.Printf("Serviço encontrado: %+v", service)

	if err := db.Delete(&service).Error; err != nil {
		log.Printf("Erro ao deletar serviço: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar serviço"})
		return
	}

	log.Printf("Serviço deletado com sucesso: ID %d", serviceID)
	c.JSON(http.StatusOK, gin.H{"message": "Serviço deletado com sucesso"})
}
