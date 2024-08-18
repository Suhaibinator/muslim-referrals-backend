package api_objects

import "muslim-referrals-backend/database"

// ReferrerView represents the fields that the referrer will be able to see

type ReferrerViewCandidate struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	WorkExperience int    `json:"workExperience"`
	ResumeUrl      string `json:"resumeUrl"`
}

func ConvertDbCandidateToReferrerViewCandidate(dbCandidate *database.Candidate) *ReferrerViewCandidate {
	return &ReferrerViewCandidate{
		FirstName:      dbCandidate.User.FirstName,
		LastName:       dbCandidate.User.LastName,
		WorkExperience: dbCandidate.WorkExperience,
		ResumeUrl:      dbCandidate.ResumeUrl,
	}
}

type ReferrerViewReferrer struct {
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Company   GeneralViewCompany `json:"company"`
}

func ConvertDbReferrerToReferrerViewReferrer(dbReferrer *database.Referrer) *ReferrerViewReferrer {
	return &ReferrerViewReferrer{
		FirstName: dbReferrer.User.FirstName,
		LastName:  dbReferrer.User.LastName,
		Company:   *ConvertDbCompanyToGeneralViewCompany(&dbReferrer.Company),
	}
}

type ReferrerViewReferralRequest struct {
	ReferralRequestId      uint64                `json:"id"`
	Candidate              ReferrerViewCandidate `json:"candidate"`
	CompanyID              uint64                `json:"company_id"`
	Company                GeneralViewCompany    `json:"company"`
	PrimaryJobTitleSeeking string                `json:"job_title"`
	JobLinks               []string              `json:"job_links"`
	Summary                string                `json:"description"`
	Locations              []string              `json:"locations"`
	ReferralType           string                `json:"referral_type"`
	ReferrerViewReferrer   ReferrerViewReferrer  `json:"referrer"`
	Status                 string                `json:"status"`
}

func ConvertDbReferralRequestToReferrerViewReferralRequest(dbReferralRequest *database.ReferralRequest) *ReferrerViewReferralRequest {

	var jobLinks []string
	for _, jobLink := range dbReferralRequest.JobLinks {
		jobLinks = append(jobLinks, jobLink.JobLink)
	}

	var locations []string
	for _, location := range dbReferralRequest.Locations {
		locations = append(locations, location.Location)
	}

	return &ReferrerViewReferralRequest{
		ReferralRequestId:      dbReferralRequest.ReferralRequestId,
		Candidate:              *ConvertDbCandidateToReferrerViewCandidate(&dbReferralRequest.Candidate),
		CompanyID:              dbReferralRequest.CompanyID,
		Company:                *ConvertDbCompanyToGeneralViewCompany(&dbReferralRequest.Company),
		PrimaryJobTitleSeeking: dbReferralRequest.PrimaryJobTitleSeeking,
		JobLinks:               jobLinks,
		Summary:                dbReferralRequest.Summary,
		Locations:              locations,
		ReferralType:           string(dbReferralRequest.ReferralType),
		ReferrerViewReferrer:   *ConvertDbReferrerToReferrerViewReferrer(dbReferralRequest.Referrer),
		Status:                 string(dbReferralRequest.Status),
	}
}
