package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/schemas"
)

func SendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, gin.H{
		"Message":   msg,
		"errorCide": code,
	})
}

func SendSucces(ctx *gin.Context, op string, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Operartion from handler %s successfull", op),
		"data":    data,
	})
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type CreateSupplierResponse struct {
	Message string                   `json:"message"`
	Data    schemas.SupplierResponse `json:"data"`
}

type DeleteSupplierResponse struct {
	Message string                   `json:"message"`
	Data    schemas.SupplierResponse `json:"data"`
}
type ShowSupplierResponse struct {
	Message string                   `json:"message"`
	Data    schemas.SupplierResponse `json:"data"`
}
type ListSuppliersResponse struct {
	Message string                     `json:"message"`
	Data    []schemas.SupplierResponse `json:"data"`
}
type UpdateSupplierResponse struct {
	Message string                   `json:"message"`
	Data    schemas.SupplierResponse `json:"data"`
}

type CreateSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

type DeleteSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

type ShowSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

type ListSupplierServicesResponse struct {
	Message string                     `json:"message"`
	Data    []schemas.SupplierServiceResponse `json:"data"`
}

type UpdateSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

type CreateSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

type DeleteSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

type ShowSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

type ListSupplierCategoriesResponse struct {
	Message string                     `json:"message"`
	Data    []schemas.SupplierCategoryResponse `json:"data"`
}

type UpdateSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}
