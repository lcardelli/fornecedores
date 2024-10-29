package handler

// Create Supplier Product
type CreateSupplierProductRequest struct {
	Name       string `json:"name"`
	SupplierID uint   `json:"supplier_id"`
	ProductID  uint   `json:"product_id"`
}

// Add tag to specify the table name
func (CreateSupplierProductRequest) TableName() string {
	return "supplier_products" // Nome da tabela correta
}

// Validate the request
func (r *CreateSupplierProductRequest) Validate() error {
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}
	if r.SupplierID == 0 {
		return errParamIsRequired("supplier_id", "uint")
	}
	if r.ProductID == 0 {
		return errParamIsRequired("product_id", "uint")
	}
	return nil
}
