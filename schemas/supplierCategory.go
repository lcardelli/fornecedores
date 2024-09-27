package schemas

import (
	"gorm.io/gorm"
)

type SupplierCategory struct {
	gorm.Model
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}