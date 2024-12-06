package schemas

import "gorm.io/gorm"

type ContractCondicaoRescisao struct {
	gorm.Model
	Name string `json:"name"`
}

type ContractCondicaoRescisaoResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

func (ContractCondicaoRescisao) TableName() string {
	return "contract_condicoes_rescisoes"
}

