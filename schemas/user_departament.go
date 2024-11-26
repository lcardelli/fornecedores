package schemas

import (
	"gorm.io/gorm"
)

// UserDepartment representa as permissões de acesso por departamento
type UserDepartment struct {
	gorm.Model
	ID            uint   `json:"id" gorm:"primaryKey"`
	UserID        uint   `json:"user_id"`
	Department    string `json:"department"`
	ViewSuppliers bool   `json:"view_suppliers"`
	ViewLicenses  bool   `json:"view_licenses"`
	AdminSuppliers bool  `json:"admin_suppliers"`
	AdminLicenses  bool  `json:"admin_licenses"`
}

