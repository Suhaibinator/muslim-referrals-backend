package models

import "time"

type Resume struct {
	Id        uint       `gorm:"primary_key" json:"id"`
	FileData  []byte     `json:"fileData"`
	UserId    uint       `json:"userId"`
	User      User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
