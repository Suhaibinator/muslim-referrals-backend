package api

type UserProfile struct {
	ID              string `json:"id"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	WorkExperience  string `json:"work_experience,omitempty"`
	CurrentCompany  string `json:"current_company,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
	CurrentTitle    string `json:"current_title,omitempty"`
	LinkedInLink    string `json:"linkedin_link,omitempty"`
	Website         string `json:"website,omitempty"`
}

type Referrer struct {
	UserID         string `json:"id"` // This ties the Referrer to a User Profile ID
	CorporateEmail string `json:"corporate_email"`
	Company        string `json:"company"`
}

type Requester struct {
	UserID             string   `json:"id"` // This ties the Requester to a User Profile ID
	ReferralRequestIDs []string `json:"referral_request_ids"`
	ResumeID           string   `json:"resume_id"`
}

type Company struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	CorporateEmailDomain []string `json:"corporate_email_domain"`
}

type ReferralRequest struct {
	ID           string   `json:"id"`
	ReferralType string   `json:"referral_type"`
	JobTitle     string   `json:"job_title"`
	CompanyID    string   `json:"company_id"`
	Locations    []string `json:"locations"`
	Description  string   `json:"description"`
	JobLinks     []string `json:"job_links,omitempty"`
	Status       string   `json:"status"`
	RequesterID  string   `json:"requester_id"`
	ReferrerID   string   `json:"referrer_id,omitempty"`
}

type Resume struct {
	ID      string `json:"resume_id"`
	Content []byte // Actual content will be stored in the database or file storage
}
