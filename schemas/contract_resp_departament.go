package schemas

import "gorm.io/gorm"

type ContractRespDepartament struct {
	gorm.Model
	Email string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
	DepartmentID uint `json:"department_id" gorm:"not null;index"`
	Department ContractDepartament `json:"department" gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type ContractRespDepartamentResponse struct {
	ID uint `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
}

func (ContractRespDepartament) TableName() string {
	return "contract_resp_departaments"
}

