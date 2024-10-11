package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/lcardelli/fornecedores/schemas"
    "net/http"
    "strconv"
)

func GetServicesByCategoryHandler(c *gin.Context) {
    categoryID := c.Query("category")

    var services []schemas.Service
    if categoryID != "" {
        categoryIDInt, _ := strconv.Atoi(categoryID)
        if err := db.Joins("JOIN supplier_services ON services.id = supplier_services.service_id").
            Joins("JOIN suppliers ON suppliers.id = supplier_services.supplier_id").
            Where("suppliers.category_id = ?", categoryIDInt).
            Distinct().Find(&services).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
            return
        }
    } else {
        if err := db.Find(&services).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços"})
            return
        }
    }

    c.JSON(http.StatusOK, services)
}