package schemas

import (
	"gorm.io/gorm"
)

type SupplierCategory struct {
	gorm.Model
	Name string `json:"name"`

	
}