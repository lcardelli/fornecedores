package schemas

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string `json:"name" gorm:"size:255;not null"`
	ServiceID  uint   `json:"service_id" gorm:"not null"`
	Service    Service `json:"service" gorm:"foreignKey:ServiceID"`
}

type ProductResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ServiceID  uint   `json:"service_id"`
	Service    ServiceResponse `json:"service"`
}
