package schemas

import (
	"gorm.io/gorm"
)

// UserDepartment representa as permissões de acesso por departamento
type UserDepartment struct {
	gorm.Model
	UserID        uint   `json:"user_id"`
	User          *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Department    string `json:"department"`
	ViewSuppliers bool   `json:"view_suppliers"`
	ViewLicenses  bool   `json:"view_licenses"`
}

