package api

import (
	"fmt"
	"log"
	"muslim-referrals-backend/database"
	"muslim-referrals-backend/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	Router   *mux.Router
	dbDriver *database.DbDriver
	service  *service.Service
}

// CORS middleware function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's an OPTIONS request, return 200
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func NewHttpServer(service *service.Service, dbd *database.DbDriver) *HttpServer {
	router := mux.NewRouter() // Create a new mux Router

	router.Use(corsMiddleware) // Use the CORS middleware

	httpServer := &HttpServer{
		Router:   router,
		dbDriver: dbd,
		service:  service,
	}

	httpServer.SetupRoutes() // Setup routes with handlers that have access to the DbDriver

	return httpServer
}

// Logging middleware function
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s | %s", r.Method, r.URL.Path, r.Header.Get("CF-Connecting-IP"))
		next.ServeHTTP(w, r)
		log.Printf("Completed [%s] %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func (hs *HttpServer) setupUserRoutes(r *mux.Router) {
	r.Use(loggingMiddleware)

	r.HandleFunc("/user/update", hs.UserUpdateUserHandler).Methods("PUT")

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

func (hs *HttpServer) setupLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login", hs.LoginHandler).Methods("GET")
}

func (hs *HttpServer) SetupRoutes() {

	// Set up your API routes
	apiRouter := hs.Router.PathPrefix("/api").Subrouter()
	hs.setupUserRoutes(apiRouter)
	hs.setupCandidateRoutes(apiRouter)
	hs.setupReferrerRoutes(apiRouter)

	// Set up the login route
	hs.setupLoginRoutes(hs.Router)

	// Serve static files from ./frontend_build for all other routes
	staticFileDirectory := http.Dir("./frontend_build")
	staticFileHandler := http.FileServer(staticFileDirectory)

	// This will serve files under ./frontend_build for any route not matched above
	hs.Router.PathPrefix("/").Handler(staticFileHandler)
}

func (hs *HttpServer) GetUserIDFromContext(r *http.Request) (uint64, error) {
	// Get the user_id from the context
	authToken, err := r.Cookie("auth")
	if err != nil || authToken == nil {
		log.Printf("[GetUserIDFromContext] Error getting auth token: %v\n", err)
		return 0, err
	}
	authTokenValue := authToken.Value
	userId, _, err := hs.service.GetUserIdFromTokenDigest(r.Context(), authTokenValue)
	return userId, err
}

func (hs *HttpServer) StartServer(port string) {
	// Start the server on port
	log.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), hs.Router)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
