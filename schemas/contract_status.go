package schemas

import "gorm.io/gorm"

type ContractStatus struct {
	gorm.Model
	Name string `json:"name"`
}

type ContractStatusResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (ContractStatus) TableName() string {
	return "contract_statuses"
}
