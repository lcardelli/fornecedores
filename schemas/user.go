package schemas

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"unique"`
	Avatar string `json:"avatar"`
	Admin  bool   `json:"admin" gorm:"default:false"`
}
