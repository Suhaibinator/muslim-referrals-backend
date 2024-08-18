package api

import (
	"muslim-referrals-backend/database"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	Router   *mux.Router
	dbDriver *database.DbDriver
}

func NewHttpServer(dbd *database.DbDriver) *HttpServer {
	router := mux.NewRouter() // Create a new mux Router

	httpServer := &HttpServer{
		Router:   router,
		dbDriver: dbd,
	}

	httpServer.SetupRoutes() // Setup routes with handlers that have access to the DbDriver

	return httpServer
}

func (hs *HttpServer) setupUserRoutes(r *mux.Router) {
	r.HandleFunc("/user/create", hs.UserCreateUserHandler).Methods("POST")

	// For all these requests, we have access to the user_id
	r.HandleFunc("/user", hs.UserGetUserHandler).Methods("GET")
	r.HandleFunc("/user/company/create", hs.UserCreateCompanyHandler).Methods("POST")
	r.HandleFunc("/user/company/get/all", hs.UserGetAllCompaniesHandler).Methods("GET")
	r.HandleFunc("/user/company/get/{company_id}", hs.UserGetCompanyHandler).Methods("GET")

	// Get my own data
	r.HandleFunc("/user/referrer/create", hs.UserCreateReferrerHandler).Methods("POST")
	r.HandleFunc("/user/referrer/update", hs.UserUpdateReferrerHandler).Methods("PUT")
	r.HandleFunc("/user/referrer/get", hs.UserGetReferrerHandler).Methods("GET")
	r.HandleFunc("/user/referrer/delete", hs.UserDeleteReferrerHandler).Methods("DELETE")

	r.HandleFunc("/user/candidate/create", hs.UserCreateCandidateHandler).Methods("POST")
	r.HandleFunc("/user/candidate/update", hs.UserUpdateCandidateHandler).Methods("PUT")
	r.HandleFunc("/user/candidate/get", hs.UserGetCandidateHandler).Methods("GET")
	r.HandleFunc("/user/candidate/delete", hs.UserDeleteCandidateHandler).Methods("DELETE")
}

func (hs *HttpServer) setupReferrerRoutes(r *mux.Router) {

	// For all these requests, we have access to the referrer_id
	r.HandleFunc("/referrer/referral_requests/all", hs.ReferrerGetAllReferralRequestsHandler).Methods("GET")
	r.HandleFunc("/referrer/referral_requests/company/{company_id}", hs.ReferrerGetReferralRequestsByCompanyHandler).Methods("GET")
	r.HandleFunc("/referrer/referral_requests/{request_id}", hs.ReferrerGetReferralRequestHandler).Methods("GET")

	// TODO: Implement this, discuss with PM
	// r.HandleFunc("/referrer/refer/{referral_request_id}", hs.ReferrerDeleteReferral).Methods("DELETE")
	// r.HandleFunc("/referrer/refer/{referral_request_id}", hs.ReferrerCreateReferral).Methods("POST")
}

func (hs *HttpServer) setupCandidateRoutes(r *mux.Router) {

	// For all these requests, we have access to the candidate_id
	r.HandleFunc("/candidate/referral_request/create", hs.CandidateCreateReferralRequestHandler).Methods("POST")
	r.HandleFunc("/candidate/referral_request/update", hs.CandidateUpdateReferralRequestHandler).Methods("PUT")
	r.HandleFunc("/candidate/referral_request/delete/{referral_request_id}", hs.CandidateDeleteReferralRequestHandler).Methods("DELETE")

	r.HandleFunc("/candidate/referral_request/get/all", hs.CandidateGetAllReferralRequestsHandler).Methods("GET")
	r.HandleFunc("/candidate/referral_request/get/{referral_request_id}", hs.CandidateGetReferralRequestHandler).Methods("GET")
}

func (hs *HttpServer) SetupRoutes() {
	hs.setupUserRoutes(hs.Router)
	hs.setupCandidateRoutes(hs.Router)
	hs.setupReferrerRoutes(hs.Router)
}

func (hs *HttpServer) GetUserIDFromContext(r *http.Request) (uint64, error) {
	// Get the user_id from the context
	//userID := r.Context().Value("user_id").(uint64)
	//return userID
	// TODO: Implement this
	return 1, nil
}

func (hs *HttpServer) StartServer() {
	// Start the server on port 8080
	http.ListenAndServe(":8080", hs.Router)
}
