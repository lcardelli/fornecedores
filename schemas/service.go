package schemas

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name       string `json:"name"`
	CategoryID uint   `json:"category_id"`
	Category   SupplierCategory `json:"category" gorm:"foreignKey:CategoryID"`
}

type ServiceResponse struct {
	ID         uint                    `json:"id"`
	Name       string                  `json:"name"`
	CategoryID uint                    `json:"category_id"`
	Category   SupplierCategoryResponse `json:"category"`
}
