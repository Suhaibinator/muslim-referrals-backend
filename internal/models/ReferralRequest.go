package models

import "time"

type ReferralType string

const (
	Internship ReferralType = "Internship"
	FullTime   ReferralType = "Full-Time"
	PartTime   ReferralType = "Part-Time"
	Contract   ReferralType = "Contract"
)

type Status string

const (
	ReferralRequested Status = "Referral Requested"
	ReferredForJob    Status = "Referred for Job"
	JobOfferReceived  Status = "Job Offer Received"
	JobOfferAccepted  Status = "Job Offer Accepted"
	JobOfferDeclined  Status = "Job Offer Declined"
)

type ReferralRequest struct {
	Id            uint         `gorm:"primary_key" json:"id"`
	JobTitle      string       `gorm:"not null" json:"jobTitle"`
	JobLinks      []string     `gorm:"not null" json:"jobLinks"`
	Description   string       `gorm:"not null" json:"description"`
	Location      string       `json:"location"`
	ReferralType  ReferralType `gorm:"not null" json:"referralType"`
	ReferrerId    *uint        `json:"referrerId"`
	Referrer      *Referrer    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"referrer"`
	CompanyId     uint         `json:"companyId"`
	Company       Company      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"company"`
	currentStatus Status       `gorm:"not null" json:"currentStatus"`
	CreatedAt     time.Time    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time    `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt     *time.Time   `json:"deletedAt"`
}
