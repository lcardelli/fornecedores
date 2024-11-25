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
	SoftwareID    uint      `json:"software_id"`
	Software      Software  `json:"software,omitempty"`
	LicenseKey    string    `json:"license_key"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Type          string    `json:"type"`
	PurchaseDate  time.Time `json:"purchase_date"`
	ExpiryDate    time.Time `json:"expiry_date"`
	Quantity      int       `json:"quantity"`
	UsedQuantity  int       `json:"used_quantity"`
	Seats         int       `json:"seats"`
	Department    string    `json:"department"`
	Cost          float64   `json:"cost"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
	AssignedUsers []User    `json:"assigned_users,omitempty" gorm:"many2many:license_users;"`
}

type LicenseUser struct {
	LicenseID uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	AssignedAt time.Time
}

func (l *License) CalculateStatus() {
	now := time.Now()
	
	// Se não tem data de expiração, considera como perpétua
	if l.ExpiryDate.IsZero() {
		l.Status = "Ativa"
		return
	}

	// Calcula a diferença em dias até a expiração
	daysUntilExpiry := l.ExpiryDate.Sub(now).Hours() / 24

	switch {
	case now.After(l.ExpiryDate):
		l.Status = "Vencida"
	case daysUntilExpiry <= 30: // Se faltam 30 dias ou menos
		l.Status = "Próxima ao vencimento"
	default:
		l.Status = "Ativa"
	}
}

func (l *License) BeforeSave(tx *gorm.DB) error {
	l.CalculateStatus()
	return nil
}

func (l *License) AfterFind(tx *gorm.DB) error {
	l.CalculateStatus()
	return nil
} 