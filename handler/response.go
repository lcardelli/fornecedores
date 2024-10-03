package handler

import "github.com/gin-gonic/gin"

func SendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, gin.H{
		"Message": msg,
		"errorCide": code,
	})
}
