package handler

import "fmt"

// Error message for required parameters
func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param: %s (type: %s) is required", name, typ)
}

// Create Supplier
type CreateSupplierRequest struct {
	Name       string                         `json:"name"`
	CNPJ       string                         `json:"cnpj"`
	Email      string                         `json:"email"`
	Phone      string                         `json:"phone"`
	Address    string                         `json:"address"`
	CategoryID uint                           `json:"category_id"`
	Services   []CreateSupplierServiceRequest `json:"services" gorm:"foreignKey:SupplierID"`
}

// Adicione a tag para especificar o nome da tabela
func (CreateSupplierRequest) TableName() string {
	return "suppliers" // Nome da tabela correta
}

// Validate the request
func (r *CreateSupplierRequest) Validate() error {
	if r.Name == "" && r.CNPJ == "" && r.Email == "" && r.Phone == "" && r.Address == "" && r.CategoryID == 0 && len(r.Services) == 0 {
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
	if len(r.Services) == 0 {
		return errParamIsRequired("services", "array")
	}
	return nil
}

// Create Supplier Service
type CreateSupplierServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SupplierID  uint    `json:"supplier_id"` // Não defina o ID aqui
}

// Add tag to specify the table name
func (CreateSupplierServiceRequest) TableName() string {
	return "supplier_services" // Nome da tabela correta
}

// Add tag to specify the table name
func (r *CreateSupplierServiceRequest) Validate() error {
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}
	if r.Description == "" {
		return errParamIsRequired("description", "string")
	}
	if r.Price == 0 {
		return errParamIsRequired("price", "float64")
	}
	return nil
}

// Create Supplier Category
type CreateSupplierCategoryRequest struct {
	Name string `json:"name"`
}

// Add tag to specify the table name
func (CreateSupplierCategoryRequest) TableName() string {
	return "supplier_categories" // Nome da tabela correta
}

// Validate the request
func (r *CreateSupplierCategoryRequest) Validate() error {
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}
	return nil
}

// Update Supplier
type UpdateSupplierRequest struct {
	Name       string                         `json:"name"`
	CNPJ       string                         `json:"cnpj"`
	Email      string                         `json:"email"`
	Phone      string                         `json:"phone"`
	Address    string                         `json:"address"`
	CategoryID uint                           `json:"category_id"`
	Services   []CreateSupplierServiceRequest `json:"services" gorm:"foreignKey:SupplierID"`
}

// Validate the request
func (r *UpdateSupplierRequest) Validate() error {
	// Verifica se todos os campos estão vazios
	if r.Name == "" && r.CNPJ == "" && r.Email == "" && r.Phone == "" && r.Address == "" && r.CategoryID == 0 && len(r.Services) == 0 {
		return fmt.Errorf("request body is empty") // Retorna erro se todos os campos estiverem vazios
	}
	// Verifica se pelo menos um campo foi fornecido
	if r.Name == "" && r.CNPJ == "" && r.Email == "" && r.Phone == "" && r.Address == "" && r.CategoryID == 0 && len(r.Services) == 0 {
		return fmt.Errorf("at least one field on request field must be provided")
	}
	return nil
}

// Update Supplier Service
type UpdateSupplierServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SupplierID  uint    `json:"supplier_id"`
}

// Validate the request
func (r *UpdateSupplierServiceRequest) Validate() error {
	if r.Name == "" || r.Description == "" || r.Price == 0 || r.SupplierID == 0 {
		return nil
	}
	return fmt.Errorf("at least one field on request field must be provided")
}

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
