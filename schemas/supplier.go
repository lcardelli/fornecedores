package schemas

import (
	"gorm.io/gorm")

type Supplier struct {
	gorm.Model
    Name        string     `json:"name"`
    CNPJ        string     `json:"cnpj" gorm:"unique"`    
    Email       string     `json:"email"`
    Phone       string     `json:"phone"`
    Address     string     `json:"address"`
    CategoryID  uint        `json:"category_id"`                                      
    Category   SupplierCategory  `json:"category" gorm:"foreignKey:CategoryID"`        
    Services    []SupplierService  `json:"services" gorm:"foreignKey:SupplierID"`
}


