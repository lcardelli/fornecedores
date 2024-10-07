package schemas

import (
	"time"

	"gorm.io/gorm"
)

type SupplierCategory struct {
	gorm.Model
	Name string `json:"name"`
}

type SupplierCategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}