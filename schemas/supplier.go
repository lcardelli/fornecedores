package schemas

import (
	"time"

	"gorm.io/gorm"
)

// SupplierLink representa o vínculo entre um fornecedor externo e as categorias/serviços locais
type SupplierLink struct {
	ID         uint           `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	CNPJ       string         `gorm:"type:varchar(18);uniqueIndex"` // Alterado aqui
	CategoryID uint
	Category   SupplierCategory
	Services   []SupplierService
}


type SupplierLinkResponse struct {
	ID                uint                `json:"id"`
	CNPJ              string              `json:"cnpj"`
	CategoryID        uint                `json:"category_id"`
	Category          SupplierCategory    `json:"category"`
	Services          []ServiceResponse   `json:"services"`
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
