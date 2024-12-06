package schemas

import "gorm.io/gorm"

type ContractDepartament struct {
	gorm.Model
	Name string `json:"name"`
}

type ContractDepartamentResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

func (ContractDepartament) TableName() string {
	return "contract_departaments"
}

