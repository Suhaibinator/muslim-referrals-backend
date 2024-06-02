package models

import "time"

type Organization struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	IsSupported bool    `json:"is_supported"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}