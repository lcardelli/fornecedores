package handler

import "github.com/lcardelli/fornecedores/schemas"

// Create Supplier Category
type CreateSupplierCategoryResponse struct {
	Message string                           `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

// Delete Supplier Category
type DeleteSupplierCategoryResponse struct {
	Message string                           `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

// Show Supplier Category
type ShowSupplierCategoryResponse struct {
	Message string                           `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}

// List Supplier Categories
type ListSupplierCategoriesResponse struct {
	Message string                             `json:"message"`
	Data    []schemas.SupplierCategoryResponse `json:"data"`
}

// Update Supplier Category
type UpdateSupplierCategoryResponse struct {
	Message string                           `json:"message"`
	Data    schemas.SupplierCategoryResponse `json:"data"`
}
