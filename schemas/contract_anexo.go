package schemas

import "gorm.io/gorm"

type ContractAnexo struct {
	gorm.Model
	ContractID uint 
	Contract Contract `json:"contract" gorm:"foreignKey:ContractID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name string `json:"name" gorm:"not null"`
	Path string `json:"path" gorm:"not null"`
	FileType string `json:"file_type" gorm:"not null"`
	FileSize int64 `json:"file_size" gorm:"not null"`
}

type ContractAnexoResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func (ContractAnexo) TableName() string {
	return "contract_anexos"
}

