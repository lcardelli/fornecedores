package handler

import "fmt"

// Update Supplier Service
type UpdateSupplierServiceRequest struct {
	Name        string  `json:"name"`
	ServiceID   uint    `json:"service_id"`
}

// Validate the request
func (r *UpdateSupplierServiceRequest) Validate() error {
	if r.Name == "" || r.ServiceID == 0 {
		return nil
	}
	return fmt.Errorf("at least one field on request field must be provided")
}
