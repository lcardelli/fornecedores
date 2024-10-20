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

// Create Supplier
type CreateSupplierResponse struct {
	Message string                      `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// Delete Supplier
type DeleteSupplierResponse struct {
	Message string                      `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// Show Supplier
type ShowSupplierResponse struct {
	Message string                      `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// List Suppliers
type ListSuppliersResponse struct {
	Message string                        `json:"message"`
	Data    []schemas.SupplierLinkResponse `json:"data"`
}

// Update Supplier
type UpdateSupplierResponse struct {
	Message string                      `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// Create Supplier Service
type CreateSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// Delete Supplier Service
type DeleteSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// Show Supplier Service
type ShowSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// List Supplier Services
type ListSupplierServicesResponse struct {
	Message string                     `json:"message"`
	Data    []schemas.SupplierServiceResponse `json:"data"`
}

// Update Supplier Service
type UpdateSupplierServiceResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// Create Supplier Category
type CreateSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

// Delete Supplier Category
type DeleteSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

// Show Supplier Category
type ShowSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

// List Supplier Categories
type ListSupplierCategoriesResponse struct {
	Message string                     `json:"message"`
	Data    []schemas.SupplierCategoryResponse `json:"data"`
}

// Update Supplier Category
type UpdateSupplierCategoryResponse struct {
	Message string                     `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}
