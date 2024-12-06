package schemas

import (
	"time"

	"gorm.io/gorm"
)

type ContractAditivo struct {
	gorm.Model

	// Parent Contract Reference
	ContractID uint     `json:"contract_id" validate:"required" gorm:"not null;index"`
	Contract   Contract `json:"contract" gorm:"foreignKey:ContractID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Basic Information
	Name            string    `json:"name" validate:"required,min=3" gorm:"not null"`
	AmendmentNumber string    `json:"amendment_number" validate:"required" gorm:"not null;index"`
	Object          string    `json:"object" validate:"required" gorm:"not null"`
	Value           float64   `json:"value" validate:"required,gte=0" gorm:"not null"`
	Observations    string    `json:"observations"`

	// Date Information
	InitialDate    time.Time `json:"initial_date" validate:"required" gorm:"not null;index"`
	FinalDate      time.Time `json:"final_date" validate:"required,gtfield=InitialDate" gorm:"not null;index"`
	SignatureDate  time.Time `json:"signature_date" validate:"required" gorm:"not null;index"`

	// Amendment Details
	Type            string  `json:"type" validate:"required,oneof=PRAZO VALOR PRAZO_VALOR OUTROS" gorm:"not null;index"`
	ValueAdjustment float64 `json:"value_adjustment"`
	NewTotalValue   float64 `json:"new_total_value" gorm:"not null"`

	// Relationships com constraints
	StatusID     uint            `json:"status_id" validate:"required" gorm:"not null;index"`
	Status       ContractStatus  `json:"status" gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	
	CostCenterID uint                    `json:"cost_center_id" validate:"required" gorm:"not null;index"`
	CostCenter   ContractCentroCusto     `json:"cost_center" gorm:"foreignKey:CostCenterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	
	BranchID     uint            `json:"branch_id" validate:"required" gorm:"not null;index"`
	Branch       ContractFilial  `json:"branch" gorm:"foreignKey:BranchID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	
	DepartmentID uint                `json:"department_id" validate:"required" gorm:"not null;index"`
	Department   ContractDepartament `json:"department" gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Attachments
	Attachments []ContractAditivoAnexo `json:"attachments" gorm:"foreignKey:ContractAditivoID;constraint:OnDelete:CASCADE"`

	// Audit Fields
	CreatedBy    uint      `json:"created_by" gorm:"not null;index"`
	UpdatedBy    uint      `json:"updated_by" gorm:"index"`
	LastModified time.Time `json:"last_modified" gorm:"autoUpdateTime;index"`
}

type ContractAditivoResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Basic Information
	Name            string    `json:"name"`
	AmendmentNumber string    `json:"amendment_number"`
	Object          string    `json:"object"`
	Value           float64   `json:"value"`
	Observations    string    `json:"observations"`

	// Date Information
	InitialDate    time.Time `json:"initial_date"`
	FinalDate      time.Time `json:"final_date"`
	SignatureDate  time.Time `json:"signature_date"`

	// Amendment Details
	Type            string  `json:"type"`
	ValueAdjustment float64 `json:"value_adjustment"`
	NewTotalValue   float64 `json:"new_total_value"`

	// Relationships
	Status         ContractStatusResponse       `json:"status"`
	CostCenter     ContractCentroCustoResponse `json:"cost_center"`
	Branch         ContractFilialResponse      `json:"branch"`
	Department     ContractDepartamentResponse `json:"department"`
	
	// Contract Reference
	ContractID     uint             `json:"contract_id"`
	ContractInfo   ContractResponse `json:"contract_info"`

	// Attachments
	Attachments []ContractAditivoAnexoResponse `json:"attachments"`

	// Audit Information
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
	LastModified time.Time `json:"last_modified"`
}

type ContractAditivoAnexo struct {
	gorm.Model
	ContractAditivoID uint   `json:"contract_aditivo_id"`
	FileName          string `json:"file_name" validate:"required"`
	FileType          string `json:"file_type" validate:"required"`
	FileSize          int64  `json:"file_size" validate:"required"`
	FilePath          string `json:"file_path" validate:"required"`
}

type ContractAditivoAnexoResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	FileURL   string    `json:"file_url"`
}
