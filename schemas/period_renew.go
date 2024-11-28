package schemas

import "gorm.io/gorm"

type PeriodRenew struct {
	gorm.Model
	Name string `json:"name"`
}