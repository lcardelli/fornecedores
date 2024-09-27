package schemas

import (
	"gorm.io/gorm"
)

type SupplierCategory struct {
	gorm.Model
	ID   uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	
}