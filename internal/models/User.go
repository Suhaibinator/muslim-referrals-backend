package models

import "time"

type User struct {
	Id          uint       `gorm:"primary_key" json:"id"`
	FullName    string     `gorm:"not null" json:"fullName"`
	Email       string     `gorm:"not null" json:"email"`
	PhoneNumber string     `json:"phoneNumber"`
	ResumeID    uint       `json:"resumeId"`
	LinkedIn    *string    `json:"linkedIn,omitempty"`
	Github      *string    `json:"github,omitempty"`
	Website     *string    `json:"website,omitempty"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}
