package handler

import (
	"github.com/lcardelli/fornecedores/schemas"
)

// Departamentos disponíveis
const (
	DepartmentTI      = "TI"
	DepartmentCompras = "Compras"
	DepartmentGeral   = "Geral"
)

// GetUserDepartments retorna as permissões do usuário
func GetUserDepartments(u *schemas.User) []schemas.UserDepartment {
	var departments []schemas.UserDepartment
	db.Where("user_id = ?", u.ID).Find(&departments)
	return departments
}

// HasSupplierAccess verifica se o usuário tem acesso à área de fornecedores
func HasSupplierAccess(u *schemas.User) bool {
	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND view_suppliers = ?", u.ID, true).First(&department)
	return result.Error == nil
}

// HasLicenseAccess verifica se o usuário tem acesso à área de licenças
func HasLicenseAccess(u *schemas.User) bool {
	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND view_licenses = ?", u.ID, true).First(&department)
	return result.Error == nil
}

// HasFullAccess verifica se o usuário tem acesso a todas as áreas
func HasFullAccess(u *schemas.User) bool {
	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND view_suppliers = ? AND view_licenses = ?",
		u.ID, true, true).First(&department)
	return result.Error == nil
}

// HasSupplierAdminAccess verifica se o usuário tem acesso administrativo à área de fornecedores
func HasSupplierAdminAccess(u *schemas.User) bool {
	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND admin_suppliers = ?", u.ID, true).First(&department)
	return result.Error == nil
}

// HasLicenseAdminAccess verifica se o usuário tem acesso administrativo à área de licenças
func HasLicenseAdminAccess(u *schemas.User) bool {
	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND admin_licenses = ?", u.ID, true).First(&department)
	return result.Error == nil
}

// SetDepartmentAccess define as permissões de acesso do usuário
func SetDepartmentAccess(u *schemas.User, department string, viewSuppliers, viewLicenses, adminSuppliers, adminLicenses bool) error {
	userDepartment := schemas.UserDepartment{
		UserID:         u.ID,
		Department:     department,
		ViewSuppliers:  viewSuppliers,
		ViewLicenses:   viewLicenses,
		AdminSuppliers: adminSuppliers,
		AdminLicenses:  adminLicenses,
	}

	result := db.Where("user_id = ? AND department = ?", u.ID, department).
		Assign(schemas.UserDepartment{
			ViewSuppliers:  viewSuppliers,
			ViewLicenses:   viewLicenses,
			AdminSuppliers: adminSuppliers,
			AdminLicenses:  adminLicenses,
		}).
		FirstOrCreate(&userDepartment)

	return result.Error
}
