package models

import "time"

type Referrer struct {
	Id             uint       `gorm:"primary_key" json:"id"`
	CorporateEmail string     `gorm:"not null" json:"corporateEmail"`
	CompanyId      uint       `json:"companyId"`
	Company        Company    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"company"`
	UserId         uint       `json:"userId"`
	User           User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user"`
	CreatedAt      time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
}
