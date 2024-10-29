package handler

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
