package handler

import "fmt"

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

func (r *CreateSupplierRequest) Validate() error {
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

func (CreateSupplierServiceRequest) TableName() string {
	return "supplier_services" // Nome da tabela correta
}

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

// Correct name for table
func (CreateSupplierCategoryRequest) TableName() string {
	return "supplier_categories" // Nome da tabela correta
}

func (r *CreateSupplierCategoryRequest) Validate() error {
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}
	return nil
}
