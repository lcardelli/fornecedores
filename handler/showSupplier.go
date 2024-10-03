package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show Supplier Handler
func ShowSupplierHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "GET Supplier",
	})
}
