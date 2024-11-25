package schemas

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"` // Este campo pode armazenar a URL da foto do perfil
	Admin  bool   `json:"admin"`
}
