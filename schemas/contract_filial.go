package schemas

import "gorm.io/gorm"

type ContractFilial struct {
	gorm.Model
	Name string `json:"name"`
}

type ContractFilialResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

func (ContractFilial) TableName() string {
	return "contract_filials"
}

