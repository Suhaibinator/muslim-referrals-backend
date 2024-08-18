package database

import (
	"time"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	_ "gorm.io/gorm"
)

type Company struct {
	Id            uint64                     `gorm:"primary_key;autoIncrement" json:"id"`
	Name          string                     `gorm:"not null" json:"name"`
	Domains       []CompanyDomainAssociation `gorm:"not null" json:"domains"`
	IsSupported   bool                       `gorm:"not null" json:"is_supported"`
	AddedByUserId uint64                     `gorm:"not null" json:"added_by_user_id"`
	User          User                       `gorm:"foreignKey:AddedByUserId;references:Id"`
	CreatedAt     time.Time                  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time                  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time                 `json:"deleted_at"`
}

type CompanyDomainAssociation struct {
	CompanyId uint64 `gorm:"primaryKey;autoIncrement:false" json:"company_id"`
	Domain    string `gorm:"primaryKey;autoIncrement:false" json:"domain"`
}

type User struct {
	Id          uint64     `gorm:"primary_key;autoIncrement" json:"id"`
	FirstName   string     `gorm:"not null" json:"firstName" validate:"required"`
	LastName    string     `gorm:"not null" json:"lastName" validate:"required"`
	Email       string     `gorm:"not null;uniqueIndex" json:"email" validate:"required,email"`
	PhoneNumber string     `json:"phoneNumber" validate:"omitempty,e164"`
	PhoneExt    string     `json:"phoneExt" validate:"omitempty,numeric"`
	LinkedIn    *string    `json:"linkedIn,omitempty" validate:"omitempty,url"`
	Github      *string    `json:"github,omitempty" validate:"omitempty,url"`
	Website     *string    `json:"website,omitempty" validate:"omitempty,url"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type Candidate struct {
	CandidateId    uint64     `gorm:"primary_key;autoIncrement" json:"id"`
	UserId         uint64     `gorm:"not null;uniqueIndex;index;constraint:OnDelete:CASCADE;foreignKey:UserId;references:Id" json:"userId"`
	User           User       `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	WorkExperience int        `gorm:"not null;" json:"workExperience"`
	ResumeUrl      string     `gorm:"not null;" json:"resumeUrl"`
	CreatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
}

type Referrer struct {
	ReferrerId     uint64     `gorm:"primary_key;autoIncrement" json:"id"`
	UserId         uint64     `gorm:"not null;uniqueIndex;index;constraint:OnDelete:CASCADE;foreignKey:UserId;references:Id" json:"userId"`
	User           User       `json:"user"`
	CompanyId      uint64     `gorm:"not null;" json:"companyId"`
	Company        Company    `json:"company"`
	CorporateEmail string     `gorm:"not null;" json:"corporateEmail"`
	CreatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
}

type ReferralType string

const (
	Internship ReferralType = "Internship"
	FullTime   ReferralType = "Full-Time"
	PartTime   ReferralType = "Part-Time"
	Contract   ReferralType = "Contract"
)

type ReferralStatus string

const (
	ReferralRequested          ReferralStatus = "Referral Requested"
	ReferralSubmissionSent     ReferralStatus = "Referred for Job"
	ReferralSubmissionAccepted ReferralStatus = "Referral Accepted"
	ReferralSubmissionRejected ReferralStatus = "Referral Rejected"
	Issue                      ReferralStatus = "Issue"
)

type ReferralRequest struct {
	ReferralRequestId      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CandidateID            uint64 `gorm:"notNull" json:"candidate_id"`
	CompanyID              uint64 `gorm:"notNull" json:"company_id"`
	PrimaryJobTitleSeeking string `gorm:"notNull" json:"job_title"`
	JobLinks               []ReferralRequestJobLinksAssociation
	Summary                string `gorm:"notNull" json:"description"`
	Locations              []ReferralRequestLocationAssociation
	ReferralType           ReferralType `gorm:"notNull" json:"referral_type"`
	ReferrerId             *uint64      `gorm:"foreignKey:ReferrerId;references:ReferrerId" json:"referrer_id"`
	Referrer               *Referrer
	Status                 ReferralStatus `gorm:"notNull" json:"status"`
	CreatedAt              time.Time      `gorm:"notNull;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"notNull;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time     `json:"deleted_at,omitempty"`
	Candidate              Candidate
	Company                Company
}

type ReferralRequestJobLinksAssociation struct {
	ReferralRequestID uint64 `gorm:"primaryKey;autoIncrement:false" json:"referral_request_id"`
	JobLink           string `gorm:"primaryKey;autoIncrement:false" json:"job_link"`
}

type ReferralRequestLocationAssociation struct {
	ReferralRequestID uint64 `gorm:"primaryKey;autoIncrement:false" json:"referral_request_id"`
	Location          string `gorm:"primaryKey;autoIncrement:false" json:"location"`
}
