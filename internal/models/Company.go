package models

import "time"

type Company struct {
	Id          uint       `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Domain      string     `gorm:"not null" json:"domain"`
	IsSupported bool       `gorm:"not null" json:"is_supported"`
	Location    string     `json:"location"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
