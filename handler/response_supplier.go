package handler

import "github.com/lcardelli/fornecedores/schemas"

// Create Supplier
type CreateSupplierResponse struct {
	Message string                       `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// Delete Supplier
type DeleteSupplierResponse struct {
	Message string                       `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// Show Supplier
type ShowSupplierResponse struct {
	Message string                       `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}

// List Suppliers
type ListSuppliersResponse struct {
	Message string                         `json:"message"`
	Data    []schemas.SupplierLinkResponse `json:"data"`
}

// Update Supplier
type UpdateSupplierResponse struct {
	Message string                       `json:"message"`
	Data    schemas.SupplierLinkResponse `json:"data"`
}
