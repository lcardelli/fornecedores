package handler

import "fmt"

// Create Supplier
type CreateSupplierRequest struct {
	Name       string `json:"name"`
	CNPJ       string `json:"cnpj"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	CategoryID uint   `json:"category_id"`
	ServiceIDs []uint `json:"service_ids"` // Mudança aqui
	ProductIDs []uint `json:"product_ids"` // Adicionado
}

// Adicione a tag para especificar o nome da tabela
func (CreateSupplierRequest) TableName() string {
	return "suppliers" // Nome da tabela correta
}

// Validate the request
func (r *CreateSupplierRequest) Validate() error {
	if r.Name == "" && r.CNPJ == "" && r.Email == "" && r.Phone == "" && r.Address == "" && r.CategoryID == 0 && len(r.ServiceIDs) == 0 && len(r.ProductIDs) == 0 {
		return fmt.Errorf("request body is empty")
	}
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}
	if r.CNPJ == "" {
		return errParamIsRequired("cnpj", "string")
	}
	if r.Email == "" {
		return errParamIsRequired("email", "string")
	}
	if r.Phone == "" {
		return errParamIsRequired("phone", "string")
	}
	if r.Address == "" {
		return errParamIsRequired("address", "string")
	}
	if r.CategoryID == 0 {
		return errParamIsRequired("category_id", "uint")
	}
	if len(r.ServiceIDs) == 0 {
		return errParamIsRequired("service_ids", "array")
	}
	if len(r.ProductIDs) == 0 {
		return errParamIsRequired("product_ids", "array")
	}
	return nil
}
