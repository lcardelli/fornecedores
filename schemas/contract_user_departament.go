package schemas

import "gorm.io/gorm"

type ContractUserDepartament struct {
	gorm.Model
	Email        string              `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	DepartmentID uint                `json:"department_id" gorm:"not null;index"`
	Department   ContractDepartament `json:"department" gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ContractUserDepartament) TableName() string {
	return "contract_user_departaments"
}