package handler

import "github.com/lcardelli/fornecedores/schemas"

// Create Supplier Service
type CreateSupplierServiceResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// Delete Supplier Service
type DeleteSupplierServiceResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// Show Supplier Service
type ShowSupplierServiceResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}

// List Supplier Services
type ListSupplierServicesResponse struct {
	Message string                            `json:"message"`
	Data    []schemas.SupplierServiceResponse `json:"data"`
}

// Update Supplier Service
type UpdateSupplierServiceResponse struct {
	Message string                          `json:"message"`
	Data    schemas.SupplierServiceResponse `json:"data"`
}
