package models

import "time"

type Referrer struct {
	Id        uint       `gorm:"primary_key" json:"id"`
	FullName  string     `gorm:"not null" json:"fullName"`
	Email     string     `gorm:"not null" json:"email"`
	CompanyId uint       `json:"companyId"`
	Company   Company    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"company"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
