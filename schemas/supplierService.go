package schemas

import (
	"gorm.io/gorm"
)

type SupplierService struct {
	gorm.Model
	SupplierLinkID uint
	ServiceID      uint
	Service        Service `gorm:"constraint:OnDelete:CASCADE;"`
}

type SupplierServiceResponse struct {
	ID         uint            `json:"id"`
	SupplierID uint            `json:"supplier_id"`
	ServiceID  uint            `json:"service_id"`
	Service    ServiceResponse `json:"service"`
}
