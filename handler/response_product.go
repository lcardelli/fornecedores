package handler

import (
	"github.com/lcardelli/fornecedores/schemas"
)

// Create Supplier Product
type CreateSupplierProductResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierProductResponse `json:"data"`
}

// Delete Supplier Product
type DeleteSupplierProductResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierProductResponse `json:"data"`
}

// Show Supplier Product
type ShowSupplierProductResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierProductResponse `json:"data"`
}

// List Supplier Products
type ListSupplierProductsResponse struct {
	Message string                            `json:"message"`
	Data    []schemas.SupplierProductResponse `json:"data"`
}

// Update Supplier Product
type UpdateSupplierProductResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierProductResponse `json:"data"`
}
