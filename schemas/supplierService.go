package schemas

import (
	"gorm.io/gorm"
)

type SupplierService struct {
	gorm.Model
	ID          int    `json:"id"`
    Nome        string `json:"nome"`
    Descricao   string `json:"descricao"`
    Preco       float64 `json:"preco"`
	
}