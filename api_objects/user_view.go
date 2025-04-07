package api_objects

import (
	"time"

	"github.com/Suhaibinator/muslim-referrals-backend/database"
)

// UserView represents the fields that the user will be able to see

type UserViewUser struct {
	Id          uint64 `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	PhoneExt    string `json:"phoneExt"`
	LinkedIn    string `json:"linkedIn"`
	Github      string `json:"github"`
	Website     string `json:"website"`
}

func ConvertUserToUserViewUser(user database.User) UserViewUser {
	var LinkedIn, Github, Website string
	if user.LinkedIn != nil {
		LinkedIn = *user.LinkedIn
	}
	if user.Github != nil {
		Github = *user.Github
	}
	if user.Website != nil {
		Website = *user.Website
	}
	return UserViewUser{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		PhoneExt:    user.PhoneExt,
		LinkedIn:    LinkedIn,
		Github:      Github,
		Website:     Website,
	}
}

func ConvertUserViewUserToUser(user UserViewUser, createdAt, updatedAt time.Time, deletedAt *time.Time) database.User {
	var LinkedIn, Github, Website *string
	if user.LinkedIn != "" {
		LinkedIn = &user.LinkedIn
	}
	if user.Github != "" {
		Github = &user.Github
	}
	if user.Website != "" {
		Website = &user.Website
	}
	return database.User{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		PhoneExt:    user.PhoneExt,
		LinkedIn:    LinkedIn,
		Github:      Github,
		Website:     Website,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   deletedAt,
	}
}

type UserViewCompany struct {
	Id      uint64   `json:"id"`
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
}

func ConvertUserViewCompanyToCompany(company UserViewCompany, userid uint64, createdAt, updatedAt time.Time, deletedAt *time.Time) database.Company {
	domains := make([]database.CompanyDomainAssociation, 0)
	for _, domain := range company.Domains {
		domains = append(domains, database.CompanyDomainAssociation{Domain: domain, CompanyId: company.Id})
	}
	return database.Company{
		Id:            company.Id,
		Name:          company.Name,
		Domains:       domains,
		AddedByUserId: userid,
		IsSupported:   false,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		DeletedAt:     deletedAt,
	}
}

func ConvertCompanyToUserViewCompany(company database.Company) UserViewCompany {
	domains := make([]string, 0)
	for _, domain := range company.Domains {
		domains = append(domains, domain.Domain)
	}
	return UserViewCompany{
		Id:      company.Id,
		Name:    company.Name,
		Domains: domains,
	}
}

type UserViewReferrer struct {
	ReferrerId     uint64 `json:"id"`
	UserId         uint64 `json:"userId"`
	CompanyId      uint64 `json:"companyId"`
	CorporateEmail string `json:"corporateEmail"`
}

func ConvertReferrerToUserViewReferrer(referrer database.Referrer) UserViewReferrer {
	return UserViewReferrer{
		ReferrerId:     referrer.ReferrerId,
		UserId:         referrer.UserId,
		CompanyId:      referrer.CompanyId,
		CorporateEmail: referrer.CorporateEmail,
	}
}

func ConvertUserViewReferrerToReferrer(referrer UserViewReferrer, userId uint64, createdAt, updatedAt time.Time, deletedAt *time.Time) database.Referrer {
	return database.Referrer{
		ReferrerId:     referrer.ReferrerId,
		UserId:         userId,
		CompanyId:      referrer.CompanyId,
		CorporateEmail: referrer.CorporateEmail,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DeletedAt:      deletedAt,
	}
}

type UserViewCandidate struct {
	CandidateId    uint64 `json:"id"`
	UserId         uint64 `json:"userId"`
	WorkExperience int    `json:"workExperience"`
	ResumeUrl      string `json:"resumeUrl"`
}

func ConvertCandidateToUserViewCandidate(candidate database.Candidate) UserViewCandidate {
	return UserViewCandidate{
		CandidateId:    candidate.CandidateId,
		UserId:         candidate.UserId,
		WorkExperience: candidate.WorkExperience,
		ResumeUrl:      candidate.ResumeUrl,
	}
}

func ConvertUserViewCandidateToCandidate(candidate UserViewCandidate, userId uint64, createdAt, updatedAt time.Time, deletedAt *time.Time) database.Candidate {
	return database.Candidate{
		CandidateId:    candidate.CandidateId,
		UserId:         userId,
		WorkExperience: candidate.WorkExperience,
		ResumeUrl:      candidate.ResumeUrl,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DeletedAt:      deletedAt,
	}
}
