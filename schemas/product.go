package schemas

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name      string  `json:"name"`
	ServiceID uint    `json:"service_id"`
	Service   Service `json:"service" gorm:"foreignKey:ServiceID"`
}

type ProductResponse struct {
	ID        uint           `json:"ID"`
	Name      string         `json:"name"`
	ServiceID uint           `json:"ServiceID"`
	Service   ServiceResponse `json:"Service"`
}
