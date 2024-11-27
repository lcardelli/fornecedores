package schemas

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Licenses    []License `json:"licenses,omitempty" gorm:"foreignKey:StatusID"`
}

// TableName especifica o nome da tabela
func (Status) TableName() string {
	return "statuses"
} 