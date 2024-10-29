package handler

// Create Supplier Service
type CreateSupplierServiceRequest struct {
	Name        string  `json:"name"`
	SupplierID  uint    `json:"supplier_id"` // Não defina o ID aqui
	ServiceID   uint    `json:"service_id"`
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
	return nil
}
