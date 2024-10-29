package schemas

import "gorm.io/gorm"

type SupplierProduct struct {
	gorm.Model
	SupplierLinkID uint        `gorm:"not null"`
	ProductID      uint        `gorm:"not null"`
	Product        Product     `gorm:"constraint:OnDelete:CASCADE;"`
}

type SupplierProductResponse struct {
	ID             uint        `json:"id"`
	Name           string      `json:"name"`
	SupplierLinkID uint        `json:"supplier_link_id"`
	ProductID      uint        `json:"product_id"`
	Product        Product     `json:"product"`
}

