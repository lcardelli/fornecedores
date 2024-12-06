package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Contract struct {
	gorm.Model
	// Basic Information
	Name           string  `json:"name" validate:"required,min=3" gorm:"not null"`
	ContractNumber string  `json:"contract_number" validate:"required" gorm:"not null;index"`
	Object         string  `json:"object" validate:"required" gorm:"not null"`
	Value          float64 `json:"value" validate:"required,gte=0" gorm:"not null"`
	Observations   string  `json:"observations"`

	// Date Information
	InitialDate time.Time `json:"initial_date" validate:"required"`
	FinalDate   time.Time `json:"final_date" validate:"required,gtfield=InitialDate"`

	// Relationships
	StatusID uint           `json:"status_id" validate:"required"`
	Status   ContractStatus `json:"status" gorm:"foreignKey:StatusID"`

	CostCenterID uint                `json:"cost_center_id" validate:"required"`
	CostCenter   ContractCentroCusto `json:"cost_center" gorm:"foreignKey:CostCenterID"`

	BranchID uint           `json:"branch_id" validate:"required"`
	Branch   ContractFilial `json:"branch" gorm:"foreignKey:BranchID"`

	DepartmentID uint                `json:"department_id" validate:"required"`
	Department   ContractDepartament `json:"department" gorm:"foreignKey:DepartmentID"`

	TerminationConditionID uint                     `json:"termination_condition_id"`
	TerminationCondition   ContractCondicaoRescisao `json:"termination_condition" gorm:"foreignKey:TerminationConditionID"`

	// Attachments
	Attachments []ContractAnexo `json:"attachments" gorm:"foreignKey:ContractID"`

	// Adicionando campos de auditoria
	CreatedBy    uint      `json:"created_by"`
	UpdatedBy    uint      `json:"updated_by"`
	LastModified time.Time `json:"last_modified" gorm:"autoUpdateTime"`
}

type ContractResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic Information
	Name           string  `json:"name"`
	ContractNumber string  `json:"contract_number"`
	Object         string  `json:"object"`
	Value          float64 `json:"value"`
	Observations   string  `json:"observations"`

	// Date Information
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`

	// Relationships
	Status               ContractStatusResponse           `json:"status"`
	CostCenter           ContractCentroCustoResponse      `json:"cost_center"`
	Branch               ContractFilialResponse           `json:"branch"`
	Department           ContractDepartamentResponse      `json:"department"`
	TerminationCondition ContractCondicaoRescisaoResponse `json:"termination_condition"`
	Attachments          []ContractAnexoResponse          `json:"attachments"`
}
