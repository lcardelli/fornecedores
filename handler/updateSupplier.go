package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Update Supplier Handler
func UpdateSupplierHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UPDATE Supplier",
	})
}
