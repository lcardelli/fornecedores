package schemas

import "gorm.io/gorm"

type ContractCentroCusto struct {
	gorm.Model
	Name string `json:"name"`
}

type ContractCentroCustoResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}


func (ContractCentroCusto) TableName() string {
	return "contract_centro_custos"
}
