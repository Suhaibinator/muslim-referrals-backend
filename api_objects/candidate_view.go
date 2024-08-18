package api_objects

import (
	"muslim-referrals-backend/database"
	"time"
)

// CandidateView represents the fields that the candidate will be able to see

type CandidateViewReferrer struct {
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Company   GeneralViewCompany `json:"company"`
}

func ConvertDbReferrerToCandidateViewReferrer(dbReferrer *database.Referrer) *CandidateViewReferrer {
	return &CandidateViewReferrer{
		FirstName: dbReferrer.User.FirstName,
		LastName:  dbReferrer.User.LastName,
		Company:   *ConvertDbCompanyToGeneralViewCompany(&dbReferrer.Company),
	}
}

type CandidateViewCandidate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func ConvertDbCandidateToCandidateViewCandidate(dbCandidate *database.Candidate) *CandidateViewCandidate {
	return &CandidateViewCandidate{
		FirstName: dbCandidate.User.FirstName,
		LastName:  dbCandidate.User.LastName,
	}
}

type CandidateViewReferralRequest struct {
	ReferralRequestId      uint64                 `json:"id"`
	Candidate              CandidateViewCandidate `json:"candidate"`
	CompanyID              uint64                 `json:"company_id"`
	Company                GeneralViewCompany     `json:"company"`
	PrimaryJobTitleSeeking string                 `json:"job_title"`
	JobLinks               []string               `json:"job_links"`
	Summary                string                 `json:"description"`
	Locations              []string               `json:"locations"`
	ReferralType           string                 `json:"referral_type"`
	ReferrerViewReferrer   CandidateViewCandidate `json:"referrer"`
	Status                 string                 `json:"status"`
}

func ConvertDbReferralRequestToCandidateViewReferralRequest(dbReferralRequest *database.ReferralRequest) *CandidateViewReferralRequest {

	var jobLinks []string
	for _, jobLink := range dbReferralRequest.JobLinks {
		jobLinks = append(jobLinks, jobLink.JobLink)
	}

	var locations []string
	for _, location := range dbReferralRequest.Locations {
		locations = append(locations, location.Location)
	}

	return &CandidateViewReferralRequest{
		ReferralRequestId:      dbReferralRequest.ReferralRequestId,
		Candidate:              *ConvertDbCandidateToCandidateViewCandidate(&dbReferralRequest.Candidate),
		CompanyID:              dbReferralRequest.CompanyID,
		Company:                *ConvertDbCompanyToGeneralViewCompany(&dbReferralRequest.Company),
		PrimaryJobTitleSeeking: dbReferralRequest.PrimaryJobTitleSeeking,
		JobLinks:               jobLinks,
		Summary:                dbReferralRequest.Summary,
		Locations:              locations,
		ReferralType:           string(dbReferralRequest.ReferralType),
		ReferrerViewReferrer:   *ConvertDbCandidateToCandidateViewCandidate(&dbReferralRequest.Candidate),
		Status:                 string(dbReferralRequest.Status),
	}
}

func ConvertCandidateViewReferralRequestToDbReferralRequest(candidateViewReferralRequest CandidateViewReferralRequest, candidateId uint64, createdAt, updatedAt time.Time, deletedAt *time.Time) database.ReferralRequest {

	var jobLinks []database.ReferralRequestJobLinksAssociation
	for _, jobLink := range candidateViewReferralRequest.JobLinks {
		jobLinks = append(jobLinks, database.ReferralRequestJobLinksAssociation{JobLink: jobLink})
	}

	var locations []database.ReferralRequestLocationAssociation
	for _, location := range candidateViewReferralRequest.Locations {
		locations = append(locations, database.ReferralRequestLocationAssociation{Location: location})
	}

	return database.ReferralRequest{
		ReferralRequestId:      candidateViewReferralRequest.ReferralRequestId,
		CandidateID:            candidateId,
		CompanyID:              candidateViewReferralRequest.CompanyID,
		PrimaryJobTitleSeeking: candidateViewReferralRequest.PrimaryJobTitleSeeking,
		JobLinks:               jobLinks,
		Summary:                candidateViewReferralRequest.Summary,
		Locations:              locations,
		ReferralType:           database.ReferralType(candidateViewReferralRequest.ReferralType),
		Status:                 database.ReferralStatus(candidateViewReferralRequest.Status),
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		DeletedAt:              deletedAt,
	}
}
