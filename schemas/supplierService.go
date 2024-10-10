package schemas

import (
	"gorm.io/gorm"
)

type SupplierService struct {
	gorm.Model
	SupplierID uint    `json:"supplier_id"`
	ServiceID  uint    `json:"service_id"`
	Service    Service `json:"service" gorm:"foreignKey:ServiceID"`
}

type SupplierServiceResponse struct {
	ID         uint            `json:"id"`
	SupplierID uint            `json:"supplier_id"`
	ServiceID  uint            `json:"service_id"`
	Service    ServiceResponse `json:"service"`
}
