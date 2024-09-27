package schemas

import (
	"gorm.io/gorm"
	"time"
)

type Supplier struct {
	gorm.Model
	ID          uint        `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string     `json:"name"`
    CNPJ        string     `json:"cnpj"`
    Email       string     `json:"email"`
    Phone       string     `json:"phone"`
    Address     string     `json:"address"`
    CategoryID  uint        `json:"category_id"`                                      
    Category   SupplierCategory  `json:"category" gorm:"foreignKey:CategoryID"`        
    Services    []SupplierService  `json:"services" gorm:"foreignKey:SupplierID"`
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
}


