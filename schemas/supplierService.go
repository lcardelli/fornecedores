package schemas

import (
	"gorm.io/gorm"
)

type SupplierService struct {
	gorm.Model
	ID          int    `json:"id"`
    Name        string `json:"name"`
    Description   string `json:"description"`
    Price       float64 `json:"price"`
	SupplierID uint `json:"supplier_id"`
}