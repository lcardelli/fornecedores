package handler

import "fmt"

// Update Supplier Category
type UpdateSupplierCategoryRequest struct {
	Name string `json:"name"`
}

// Validate the request
func (r *UpdateSupplierCategoryRequest) Validate() error {
	if r.Name == "" {
		return nil
	}
	return fmt.Errorf("at least one field on request field must be provided")
}
