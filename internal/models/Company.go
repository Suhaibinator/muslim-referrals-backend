package models

import "time"

type Company struct {
	Id          uint       `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Domain      string     `gorm:"not null" json:"domain"`
	IsSupported bool       `gorm:"not null" json:"isASpported"`
	Location    string     `json:"location"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}
