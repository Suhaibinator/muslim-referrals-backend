package api

import (
	"encoding/json"
	"github.com/Suhaibinator/muslim-referrals-backend/api_objects"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CandidateCreateReferralRequestHandler handles the creation of a referral request by a candidate
func (hs *HttpServer) CandidateCreateReferralRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called CandidateCreateReferralRequestHandler")

	// Step 1: Decode the incoming JSON request into a ReferralRequest struct
	var request api_objects.CandidateViewReferralRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: Authenticate and retrieve the candidate/user ID
	userID, userIdRetrievalErr := hs.GetUserIDFromContext(r)
	if userIdRetrievalErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Step 3: Validate that the authenticated user is the candidate making the request
	candidate := hs.dbDriver.GetCandidateByUserId(userID)
	if candidate == nil || candidate.UserId != userID {
		http.Error(w, "Candidate not found or unauthorized", http.StatusForbidden)
		return
	}

	referralRequest := api_objects.ConvertCandidateViewReferralRequestToDbReferralRequest(request, candidate.CandidateId, time.Now(), time.Now(), nil)

	// Step 4: Create the referral request in the database
	createdRequest, err := hs.dbDriver.CreateReferralRequest(&referralRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 5: Convert the domain model to the view model, if necessary, or directly marshal the created object
	response, err := json.Marshal(api_objects.ConvertDbReferralRequestToCandidateViewReferralRequest(createdRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 6: Set header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// CandidateUpdateReferralRequestHandler handles updating a referral request by a candidate
func (hs *HttpServer) CandidateUpdateReferralRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called CandidateUpdateReferralRequestHandler")

	var requestUpdate api_objects.CandidateViewReferralRequest
	if err := json.NewDecoder(r.Body).Decode(&requestUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, userIdRetrievalErr := hs.GetUserIDFromContext(r)
	if userIdRetrievalErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := hs.dbDriver.GetCandidateByUserId(userID)
	if candidate == nil || candidate.UserId != userID {
		http.Error(w, "Candidate not found or unauthorized", http.StatusForbidden)
		return
	}

	existingRequest := hs.dbDriver.GetReferralRequestById(requestUpdate.ReferralRequestId)

	if existingRequest == nil {
		http.Error(w, "Referral request not found", http.StatusInternalServerError)
		return
	}

	if existingRequest.Candidate.CandidateId != candidate.CandidateId {
		http.Error(w, "Unauthorized to update this referral request", http.StatusForbidden)
		return
	}

	updatedRequest := api_objects.ConvertCandidateViewReferralRequestToDbReferralRequest(requestUpdate, candidate.CandidateId, existingRequest.CreatedAt, time.Now(), nil)
	dbResult, dbUpdateErr := hs.dbDriver.UpdateReferralRequest(&updatedRequest)
	if dbUpdateErr != nil {
		http.Error(w, dbUpdateErr.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(api_objects.ConvertDbReferralRequestToCandidateViewReferralRequest(dbResult))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// CandidateDeleteReferralRequestHandler handles the deletion of a referral request by a candidate
func (hs *HttpServer) CandidateDeleteReferralRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called CandidateDeleteReferralRequestHandler")

	// Authenticate and retrieve the candidate/user ID
	userID, userIdRetrievalErr := hs.GetUserIDFromContext(r)
	if userIdRetrievalErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := hs.dbDriver.GetCandidateByUserId(userID)

	// Extract the referral request ID from URL parameters
	vars := mux.Vars(r)
	referralRequestID, err := strconv.ParseUint(vars["referral_request_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid referral request ID", http.StatusBadRequest)
		return
	}

	// Retrieve the referral request to ensure it belongs to the authenticated user
	referralRequest := hs.dbDriver.GetReferralRequestById(referralRequestID)
	if referralRequest == nil || candidate == nil || referralRequest.CandidateID != candidate.CandidateId {
		http.Error(w, "Referral request not found or unauthorized", http.StatusForbidden)
		return
	}

	// Delete the referral request from the database
	err = hs.dbDriver.DeleteReferralRequest(referralRequest)
	if err != nil {
		http.Error(w, "Failed to delete referral request", http.StatusInternalServerError)
		return
	}

	// Successful deletion
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// CandidateGetReferralRequestHandler handles fetching a specific referral request by its ID
func (hs *HttpServer) CandidateGetReferralRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called CandidateGetReferralRequestHandler")

	userID, userIdRetrievalErr := hs.GetUserIDFromContext(r)
	if userIdRetrievalErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := hs.dbDriver.GetCandidateByUserId(userID)
	if candidate == nil {
		http.Error(w, "Candidate not found or unauthorized", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	referralRequestID, err := strconv.ParseUint(vars["referral_request_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid referral request ID", http.StatusBadRequest)
		return
	}

	referralRequest := hs.dbDriver.GetReferralRequestByIdAndCandidateId(referralRequestID, candidate.CandidateId)
	if referralRequest == nil {
		http.Error(w, "Referral request not found or unauthorized", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(api_objects.ConvertDbReferralRequestToCandidateViewReferralRequest(referralRequest))
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// CandidateGetAllReferralRequestsHandler handles fetching all referral requests for a candidate
func (hs *HttpServer) CandidateGetAllReferralRequestsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called CandidateGetAllReferralRequestsHandler")
	userID, userIdRetrievalErr := hs.GetUserIDFromContext(r)
	if userIdRetrievalErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := hs.dbDriver.GetCandidateByUserId(userID)
	if candidate == nil {
		http.Error(w, "Candidate not found or unauthorized", http.StatusForbidden)
		return
	}

	referralRequests := hs.dbDriver.GetReferralRequestsByCandidateId(candidate.CandidateId)

	result := make([]api_objects.CandidateViewReferralRequest, 0)
	for _, referralRequest := range referralRequests {
		result = append(result, *api_objects.ConvertDbReferralRequestToCandidateViewReferralRequest(&referralRequest))
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
