package api

import "github.com/gorilla/mux"

func setupUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{user_id}", GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{user_id}", UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{user_id}", DeleteUserHandler).Methods("DELETE")
}

func setupReferrerRoutes(r *mux.Router) {
	r.HandleFunc("/referrers", CreateReferrerHandler).Methods("POST")
	r.HandleFunc("/referrers/{user_id}", GetReferrerHandler).Methods("GET")
	r.HandleFunc("/referrers/{referrer_id}", UpdateReferrerHandler).Methods("PUT")
	r.HandleFunc("/referrers/{user_id}", DeleteReferrerHandler).Methods("DELETE")
}

func setupRequesterRoutes(r *mux.Router) {
	r.HandleFunc("/requesters", CreateRequesterHandler).Methods("POST")
	r.HandleFunc("/requesters/{requester_id}", GetRequesterHandler).Methods("GET")
	r.HandleFunc("/requesters/{requester_id}", UpdateRequesterHandler).Methods("PUT")
	r.HandleFunc("/requesters/{requester_id}", DeleteRequesterHandler).Methods("DELETE")
	r.HandleFunc("/resumes", UploadResumeHandler).Methods("POST")
}

func setupCompanyRoutes(r *mux.Router) {
	r.HandleFunc("/companies", CreateCompanyHandler).Methods("POST")
	r.HandleFunc("/companies/{company_id}", GetCompanyHandler).Methods("GET")
	r.HandleFunc("/companies/{company_id}", UpdateCompanyHandler).Methods("PUT")
	r.HandleFunc("/companies/{company_id}", DeleteCompanyHandler).Methods("DELETE")
}

func setupReferralRequestRoutes(r *mux.Router) {
	r.HandleFunc("/referral-requests", CreateReferralRequestHandler).Methods("POST")
	r.HandleFunc("/referral-requests/{referral_request_id}", GetReferralRequestHandler).Methods("GET")
	r.HandleFunc("/referral-requests/{referral_request_id}", UpdateReferralRequestHandler).Methods("PUT")
	r.HandleFunc("/referral-requests/{referral_request_id}", DeleteReferralRequestHandler).Methods("DELETE")
	r.HandleFunc("/referral-requests/requester/{requester_id}", GetReferralRequestsByRequesterHandler).Methods("GET")
	r.HandleFunc("/referral-requests/referrer/{referrer_id}", GetReferralRequestsByReferrerHandler).Methods("GET")
	r.HandleFunc("/referral-requests/company/{company_id}", GetReferralRequestsByCompanyHandler).Methods("GET")
}

func SetupRoutes(r *mux.Router) {
	setupUserRoutes(r)
	setupReferrerRoutes(r)
	setupRequesterRoutes(r)
	setupCompanyRoutes(r)
	setupReferralRequestRoutes(r)
}

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	SetupRoutes(r)
	return r
}
