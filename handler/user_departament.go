package handler

import (
	"github.com/lcardelli/fornecedores/schemas"
	"fmt"
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
	// Admin global tem acesso a tudo
	if u.Admin {
		return true
	}
	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND view_suppliers = ?", u.ID, true).First(&department)
	return result.Error == nil
}

// HasLicenseAccess verifica se o usuário tem acesso à área de licenças
func HasLicenseAccess(u *schemas.User) bool {
	// Admin global tem acesso a tudo
	if u.Admin {
		return true
	}
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
	// Admin global tem acesso a tudo
	if u.Admin {
		return true
	}

	var department schemas.UserDepartment
	// Busca direta sem condição de admin_suppliers para debug
	result := db.Where("user_id = ?", u.ID).First(&department)
	
	fmt.Printf("Debug HasSupplierAdminAccess:\n")
	fmt.Printf("- Usuário ID: %d\n", u.ID)
	fmt.Printf("- Erro na consulta: %v\n", result.Error)
	fmt.Printf("- Departamento: %+v\n", department)
	fmt.Printf("- AdminSuppliers: %v\n", department.AdminSuppliers)
	
	return result.Error == nil && department.AdminSuppliers
}

// HasLicenseAdminAccess verifica se o usuário tem acesso administrativo à área de licenças
func HasLicenseAdminAccess(u *schemas.User) bool {
	// Admin global tem acesso a tudo
	if u.Admin {
		return true
	}

	var department schemas.UserDepartment
	// Remover a cláusula deleted_at IS NULL temporariamente para debug
	result := db.Unscoped().Where("user_id = ? AND admin_licenses = ?", u.ID, true).First(&department)
	
	fmt.Printf("Verificando acesso admin licenças para usuário %d: %v\n", u.ID, result.Error == nil)
	fmt.Printf("Departamento encontrado: %+v\n", department)
	
	return result.Error == nil
}

// HasContractAdminAccess verifica se o usuário tem acesso administrativo à área de contratos
func HasContractAdminAccess(u *schemas.User) bool {
	// Admin global tem acesso a tudo
	if u.Admin {
		return true
	}

	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND admin_contracts = ?", u.ID, true).First(&department)
	return result.Error == nil
}

func HasContractViewAccess(u *schemas.User) bool {
	// Admin global tem acesso a tudo
	if u.Admin {
		return true
	}

	var department schemas.UserDepartment
	result := db.Where("user_id = ? AND view_contracts = ?", u.ID, true).First(&department)
	return result.Error == nil
}




// SetDepartmentAccess define as permissões de acesso do usuário
func SetDepartmentAccess(u *schemas.User, departmentID uint, viewSuppliers, viewLicenses, adminSuppliers, adminLicenses, viewContracts, adminContracts bool) error {
	userDepartment := schemas.UserDepartment{
		UserID:         u.ID,
		DepartmentID:   departmentID,
		ViewSuppliers:  viewSuppliers,
		ViewLicenses:   viewLicenses,
		AdminSuppliers: adminSuppliers,
		AdminLicenses:  adminLicenses,
		ViewContracts:  viewContracts,
		AdminContracts: adminContracts,
	}

	result := db.Where("user_id = ? AND department_id = ?", u.ID, departmentID).
		Assign(schemas.UserDepartment{
			ViewSuppliers:  viewSuppliers,
			ViewLicenses:   viewLicenses,
			AdminSuppliers: adminSuppliers,
			AdminLicenses:  adminLicenses,
			ViewContracts:  viewContracts,
			AdminContracts: adminContracts,
		}).
		FirstOrCreate(&userDepartment)

	return result.Error
}
