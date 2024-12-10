package schemas

import (
	"gorm.io/gorm"
)

// UserDepartment representa as permissões de acesso por departamento
type UserDepartment struct {
	gorm.Model
	ID             uint `json:"id" gorm:"primaryKey"`
	UserID         uint `json:"user_id"`
	DepartmentID   uint `json:"department_id"`
	ViewSuppliers  bool `json:"view_suppliers"`
	ViewLicenses   bool `json:"view_licenses"`
	AdminSuppliers bool `json:"admin_suppliers"`
	AdminLicenses  bool `json:"admin_licenses"`
	ViewContracts  bool `json:"view_contracts"`
	AdminContracts bool `json:"admin_contracts"`
}
