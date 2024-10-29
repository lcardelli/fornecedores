package schemas

import (
	"time"

	"gorm.io/gorm"
)

// SupplierLink representa o vínculo entre um fornecedor externo e as categorias/serviços locais
type SupplierLink struct {
	gorm.Model
	CNPJ       string         `gorm:"type:varchar(18);uniqueIndex"` // Alterado aqui
	CategoryID uint
	Category   SupplierCategory
	Services   []SupplierService
	Products   []SupplierProduct
}


type SupplierLinkResponse struct {
	ID                uint                `json:"id"`
	CNPJ              string              `json:"cnpj"`
	CategoryID        uint                `json:"category_id"`
	Category          SupplierCategory    `json:"category"`
	Services          []ServiceResponse   `json:"services"`
	Products          []ProductResponse   `json:"products"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	DeletedAt         time.Time           `json:"deleted_at,omitempty"`
	ExternalSupplier  ExternalSupplier    `json:"external_supplier"`
}

// ExternalSupplier representa os fornecedores vindos do banco de dados externo
type ExternalSupplier struct {
	CODCOLIGADA   int
	CODCFO        string
	NOMEFANTASIA  string
	NOME          string
	CGCCFO        string
	RUA           string
	NUMERO        string
	COMPLEMENTO   string
	BAIRRO        string
	CIDADE        string
	CEP           string
	TELEFONE      string
	EMAIL         string
	CONTATO       string
	UF            string
	ATIVO         string
	TIPO          string
}

// Supplier representa o fornecedor
type Supplier struct {
	gorm.Model
	CNPJ       string `json:"cnpj"`
	CategoryID uint   `json:"category_id"`
	Services   []Service   `gorm:"many2many:supplier_services;"`
	Products   []Product   `gorm:"many2many:supplier_products;"`
}


