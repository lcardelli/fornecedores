package handler

import "fmt"

// Update Supplier
type UpdateSupplierRequest struct {
	Name       string                         `json:"name"`
	CNPJ       string                         `json:"cnpj"`
	Email      string                         `json:"email"`
	Phone      string                         `json:"phone"`
	Address    string                         `json:"address"`
	CategoryID uint                           `json:"category_id"`
	Services   []CreateSupplierServiceRequest `json:"services" gorm:"foreignKey:SupplierID"`
	Products   []CreateSupplierProductRequest `json:"products" gorm:"foreignKey:SupplierID"`
}

// Validate the request
func (r *UpdateSupplierRequest) Validate() error {
	// Verifica se todos os campos estão vazios
	if r.Name == "" && r.CNPJ == "" && r.Email == "" && r.Phone == "" && r.Address == "" && r.CategoryID == 0 && len(r.Services) == 0 && len(r.Products) == 0 {
		return fmt.Errorf("request body is empty") // Retorna erro se todos os campos estiverem vazios
	}
	// Verifica se pelo menos um campo foi fornecido
	if r.Name == "" && r.CNPJ == "" && r.Email == "" && r.Phone == "" && r.Address == "" && r.CategoryID == 0 && len(r.Services) == 0 && len(r.Products) == 0 {
		return fmt.Errorf("at least one field on request field must be provided")
	}
	return nil
}
