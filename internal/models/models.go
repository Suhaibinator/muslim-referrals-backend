package models

import "time"

type Company struct {
	Id          uint       `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Domain      string     `gorm:"not null" json:"domain"`
	IsSupported bool       `gorm:"not null" json:"is_supported"`
	Location    string     `json:"location"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type User struct {
	Id          uint       `gorm:"primary_key" json:"id"`
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
	Id             uint      `gorm:"primary_key;autoIncrement" json:"id"`
	UserId         uint      `gorm:"not null;uniqueIndex;constraint:OnDelete:CASCADE;" json:"userId"`
	User           User      `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	WorkExperience int       `gorm:"not null;" json:"workExperience"`
	ResumeUrl      string    `gorm:"not null;" json:"resumeUrl"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

type Referrer struct {
	Id             uint      `gorm:"primary_key;autoIncrement" json:"id"`
	UserId         uint      `gorm:"not null;uniqueIndex;constraint:OnDelete:CASCADE;" json:"userId"`
	User           User      `json:"user"`
	CompanyId      uint      `gorm:"not null;" json:"companyId"`
	Company        Company   `json:"company"`
	CorporateEmail string    `gorm:"not null;" json:"corporateEmail"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
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
	CandidateID            uint                                 `gorm:"primaryKey;autoIncrement:false" json:"candidateId"`
	CompanyID              uint                                 `gorm:"primaryKey;autoIncrement:false" json:"companyId"`
	PrimaryJobTitleSeeking string                               `gorm:"not null" json:"jobTitle"`
	JobLinks               []ReferralRequestJobLinksAssociation `gorm:"foreignKey:ReferralRequestID;references:CandidateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"jobLinks"`
	Summary                string                               `gorm:"not null" json:"description"`
	Locations              []ReferralRequestLocationAssociation `gorm:"foreignKey:ReferralRequestID;references:CandidateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"locations"`
	ReferralType           ReferralType                         `gorm:"not null" json:"referralType"`
	ReferrerId             *uint                                `json:"referrerId"`
	Referrer               *Referrer                            `gorm:"foreignKey:ReferrerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"referrer"`
	Status                 ReferralStatus                       `gorm:"not null" json:"status"`
	CreatedAt              time.Time                            `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt              time.Time                            `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt              *time.Time                           `json:"deletedAt"`
	Candidate              Candidate                            `gorm:"foreignKey:CandidateID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"candidate"`
	Company                Company                              `gorm:"foreignKey:CompanyID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company"`
}

type ReferralRequestJobLinksAssociation struct {
	ReferralRequestID uint   `gorm:"primaryKey" json:"referralRequestId"`
	JobLink           string `gorm:"primaryKey" json:"jobLink"`
}

type ReferralRequestLocationAssociation struct {
	ReferralRequestID uint   `gorm:"primaryKey" json:"referralRequestId"`
	Location          string `gorm:"primaryKey" json:"location"`
}
