package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete Supplier Handler
func DeleteSupplierHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "DELETE Supplier",
	})
}
