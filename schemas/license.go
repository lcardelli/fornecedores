package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Software struct {
	gorm.Model
	Name        string    `json:"name"`
	Publisher   string    `json:"publisher"`
	Description string    `json:"description"`
	Licenses    []License `json:"licenses,omitempty" gorm:"foreignKey:SoftwareID"`
}

type License struct {
	gorm.Model
	SoftwareID    uint           `json:"software_id"`
	Software      Software       `json:"software"`
	LicenseKey    string        `json:"license_key"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	Type          string        `json:"type"`
	DepartmentID  uint          `json:"department_id"`
	Department    Departament   `json:"department"`
	Cost         float64       `json:"cost"`
	SiteURL      string        `json:"site_url"`
	PurchaseDate time.Time     `json:"purchase_date"`
	ExpiryDate   time.Time     `json:"expiry_date"`
	Notes        string        `json:"notes"`
	StatusID     uint          `json:"status_id"`
	Status       Status        `json:"status"`
	PeriodRenewID *uint        `json:"period_renew_id"`
	PeriodRenew   *PeriodRenew `json:"period_renew"`
	AssignedUsers []User       `gorm:"many2many:license_users;" json:"assigned_users"`
	Blocked       bool          `json:"blocked"`
}

type LicenseUser struct {
	LicenseID  uint `gorm:"primaryKey"`
	UserID     uint `gorm:"primaryKey"`
	AssignedAt time.Time
}

func (l *License) CalculateStatus(db *gorm.DB) error {
	now := time.Now()
	var status Status

	// Se a licença estiver bloqueada, define status como cancelada
	if l.Blocked {
		if err := db.Table("statuses").Where("name = ?", "Cancelada").First(&status).Error; err != nil {
			return err
		}
	} else {
		// Se não tem data de expiração, considera como perpétua
		if l.ExpiryDate.IsZero() {
			if err := db.Table("statuses").Where("id = ?", 1).First(&status).Error; err != nil {
				return err
			}
		} else {
			// Calcula a diferença em dias até a expiração
			daysUntilExpiry := l.ExpiryDate.Sub(now).Hours() / 24

			var statusID uint
			switch {
			case now.After(l.ExpiryDate):
				statusID = 3
			case daysUntilExpiry <= 30:
				statusID = 2
			default:
				statusID = 1
			}

			if err := db.Table("statuses").Where("id = ?", statusID).First(&status).Error; err != nil {
				return err
			}
		}
	}

	l.StatusID = status.ID
	l.Status = status
	return nil
}

func (l *License) BeforeSave(tx *gorm.DB) error {
	return l.CalculateStatus(tx)
}

func (l *License) AfterFind(tx *gorm.DB) error {
	return l.CalculateStatus(tx)
}

func (l *License) BeforeCreate(tx *gorm.DB) error {
	return l.CalculateStatus(tx)
}

func (l *License) BeforeUpdate(tx *gorm.DB) error {
	return l.CalculateStatus(tx)
}
