package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// List Supplier Handler
func ListSupplierHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "LIST Supplier",
	})
}
