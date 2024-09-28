package schemas

import (
	"gorm.io/gorm"
)

type SupplierService struct {
	gorm.Model
    Name        string `json:"name"`
    Description   string `json:"description"`
    Price       float64 `json:"price"`
	SupplierID uint `json:"supplier_id"`

}