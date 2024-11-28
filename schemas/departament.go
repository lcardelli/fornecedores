package schemas

import "gorm.io/gorm"

type Departament struct {
	gorm.Model
	Name string `json:"name"`
}

func (Departament) TableName() string {
	return "departaments"
}
