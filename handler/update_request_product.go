package handler

import "fmt"

// Update Supplier Product
type UpdateSupplierProductRequest struct {
	Name       string `json:"name"`
	SupplierID uint   `json:"supplier_id"`
}

// Validate the request
func (r *UpdateSupplierProductRequest) Validate() error {
	if r.Name == "" || r.SupplierID == 0 {
		return nil
	}
	return fmt.Errorf("at least one field on request field must be provided")
}
