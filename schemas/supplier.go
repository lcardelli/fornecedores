package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Supplier struct {
	gorm.Model
	Name       string 
	CNPJ       string `gorm:"type:varchar(20);unique"` // Definindo o tipo e o comprimento
	Email      string
	Phone      string
	Address    string
	CategoryID uint
	Category   SupplierCategory `gorm:"constraint:OnDelete:CASCADE;"` // Adicionando exclusão em cascata
	Services   []SupplierService `gorm:"constraint:OnDelete:CASCADE;"` // Adicionando exclusão em cascata
}

type SupplierResponse struct {
	ID         uint                `json:"id"`
	Name       string              `json:"name"`
	CNPJ       string              `json:"cnpj"`
	Email      string              `json:"email"`
	Phone      string              `json:"phone"`
	Address    string              `json:"address"`
	CategoryID uint                `json:"category_id"`
	Category   SupplierCategory    `json:"category"`
	Services   []ServiceResponse   `json:"services"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	DeletedAt  time.Time           `json:"deleted_at,omitempty"`
}
