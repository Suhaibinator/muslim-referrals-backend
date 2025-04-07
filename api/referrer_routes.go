package api

import (
	"encoding/json"
	"github.com/Suhaibinator/muslim-referrals-backend/api_objects"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ReferrerGetAllReferralRequestsHandler handles fetching all referral requests for a referrer
func (hs *HttpServer) ReferrerGetAllReferralRequestsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called ReferrerGetAllReferralRequestsHandler")

	userID, userIdRetrievalErr := hs.GetUserIDFromContext(r)
	if userIdRetrievalErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	referrer := hs.dbDriver.GetReferrerByUserId(userID)
	if referrer == nil {
		http.Error(w, "Referrer not found or unauthorized", http.StatusForbidden)
		return
	}

	referralRequests := hs.dbDriver.GetReferralRequestsByCompanyId(referrer.CompanyId)

	result := make([]api_objects.ReferrerViewReferralRequest, 0)

	for _, referralRequest := range referralRequests {
		result = append(result, *api_objects.ConvertDbReferralRequestToReferrerViewReferralRequest(&referralRequest))
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error marshaling referral requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// ReferrerGetReferralRequestsHandler handles fetching referral requests for a referrer based on specific criteria
func (hs *HttpServer) ReferrerGetReferralRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called ReferrerGetReferralRequestsHandler")

	userID, err := hs.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	referrer := hs.dbDriver.GetReferrerByUserId(userID)
	if referrer == nil {
		http.Error(w, "Referrer not found or unauthorized", http.StatusForbidden)
		return
	}

	// Get the referral request ID from the URL
	vars := mux.Vars(r)
	referralRequestId, err := strconv.ParseUint(vars["request_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid referral request ID", http.StatusBadRequest)
		return
	}

	// Fetch referral requests based on the given status and referrer ID
	referralRequest := hs.dbDriver.GetReferralRequestById(referralRequestId)

	result := api_objects.ConvertDbReferralRequestToReferrerViewReferralRequest(referralRequest)

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error marshaling referral request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// ReferrerGetReferralRequestsByCompanyHandler handles fetching referral requests for a referrer based on the company
func (hs *HttpServer) ReferrerGetReferralRequestsByCompanyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called ReferrerGetReferralRequestsByCompanyHandler")

	userID, err := hs.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	company_id, companyIdParsingErr := strconv.ParseUint(mux.Vars(r)["company_id"], 10, 64)
	if companyIdParsingErr != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	referrer := hs.dbDriver.GetReferrerByUserId(userID)
	if referrer == nil || referrer.CompanyId != company_id {
		http.Error(w, "Referrer not found or unauthorized", http.StatusForbidden)
		return
	}

	referralRequests := hs.dbDriver.GetReferralRequestsByCompanyId(company_id)

	result := make([]api_objects.ReferrerViewReferralRequest, 0)
	for _, referralRequest := range referralRequests {
		result = append(result, *api_objects.ConvertDbReferralRequestToReferrerViewReferralRequest(&referralRequest))
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error marshaling referral requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
