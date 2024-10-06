package main

import (
	"fmt"
	"io"
	"muslim-referrals-backend/api"
	"muslim-referrals-backend/config"
	"muslim-referrals-backend/database"
	"muslim-referrals-backend/service"
	"os/signal"
	"syscall"
	"time"

	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {

	stmts, err := gormschema.New("sqlite").Load(
		&database.User{},
		&database.Company{},
		&database.Candidate{},
		&database.Referrer{},
		&database.ReferralRequest{},
		&database.ReferralRequestJobLinksAssociation{},
		&database.ReferralRequestLocationAssociation{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)

	db := database.NewDbDriver(config.DatabasePath)
	defer db.CloseDatabase()

	service := service.NewService(config.GoogleOauthConfig, db)

	httpServer := api.NewHttpServer(service, db)
	go httpServer.StartServer(config.Port)

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	fmt.Printf("Received signal: %s. Shutting down...\n", sig)

	// testCreations(db)

	// refReq := db.GetReferralRequestById(1)

	// fmt.Println("Referral Request: ", refReq)
}

func testCreations(db *database.DbDriver) {

	// Create a User
	user := database.User{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe3@email.com",
		PhoneNumber: "123-456-7890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdUser, createdUserErr := db.CreateUser(&user)
	if createdUserErr != nil {
		fmt.Println("Error creating user: ", createdUserErr)
	}

	company := database.Company{
		Name: "Tech Innovators",
		Domains: []database.CompanyDomainAssociation{{
			Domain: "techinnovators.com",
		}},
		AddedByUserId: createdUser.Id,
		IsSupported:   true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdCompany, createdCompanyErr := db.CreateCompany(&company)
	if createdCompanyErr != nil {
		fmt.Println("Error creating company: ", createdCompanyErr)
	}
	fmt.Println("Created Company: ", createdCompany)

	// Create a Candidate linked to the User
	candidate := database.Candidate{
		UserId:         createdUser.Id,
		WorkExperience: 5,
		ResumeUrl:      "http://example.com/resume.pdf",
		CreatedAt:      time.Now(),
	}

	createdCandidate, createdCandidateErr := db.CreateCandidate(&candidate)
	if createdCandidateErr != nil {
		fmt.Println("Error creating candidate: ", createdCandidateErr)
	}

	fmt.Println("Created Candidate: ", createdCandidate)

	// Create a Referrer linked to the User and Company
	referrer := database.Referrer{
		UserId:         createdUser.Id,
		CompanyId:      createdCompany.Id,
		CorporateEmail: "john.doe@techinnovators.com",
		CreatedAt:      time.Now(),
	}

	createdReferrer, createdReferrerErr := db.CreateReferrer(&referrer)
	if createdReferrerErr != nil {
		fmt.Println("Error creating referrer: ", createdReferrerErr)
	}

	// Create a ReferralRequest linked to the Candidate and Company
	referralRequest := database.ReferralRequest{
		CandidateID:            createdCandidate.CandidateId,
		CompanyID:              company.Id,
		PrimaryJobTitleSeeking: "Software Developer",
		Summary:                "Experienced developer seeking new challenges.",
		ReferralType:           database.FullTime,
		ReferrerId:             &createdReferrer.ReferrerId,
		JobLinks: []database.ReferralRequestJobLinksAssociation{
			{
				JobLink: "http://example.com/job1",
			}, {
				JobLink: "http://example.com/job2",
			}},
		Status:    database.ReferralRequested,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdReferralRequest, createdReferralRequestErr := db.CreateReferralRequest(&referralRequest)
	if createdReferralRequestErr != nil {
		fmt.Println("Error creating referral request: ", createdReferralRequestErr)
	}

	fmt.Println("Created Referral Request: ", createdReferralRequest)
}
