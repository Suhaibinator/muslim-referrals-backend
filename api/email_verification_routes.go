package api

import (
	"encoding/json"
	"errors"
	"github.com/Suhaibinator/muslim-referrals-backend/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type EmailVerificationRequestPayload struct {
	Email string `json:"email"`
}

// EmailVerificationRequestHandler handles the creation of a new email verification request.
// POST /api/email-verification
func (hs *HttpServer) EmailVerificationRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := hs.GetUserIDFromContext(r)
	if err != nil {
		log.Printf("Error getting user ID from context: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload EmailVerificationRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// TODO: Add validation to check if the email format is valid

	err = hs.service.RequestEmailVerification(userID, payload.Email)
	if err != nil {
		log.Printf("Error requesting email verification for user %d, email %s: %v", userID, payload.Email, err)
		// Don't expose internal error details directly
		http.Error(w, "Failed to process verification request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(map[string]string{"message": "Verification email request sent successfully."})
}

// EmailVerificationVerifyHandler handles the verification of an email using a code.
// GET /api/email-verification/verify/{verification_code}
func (hs *HttpServer) EmailVerificationVerifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	verificationCode, ok := vars["verification_code"]
	if !ok || verificationCode == "" {
		http.Error(w, "Verification code is missing", http.StatusBadRequest)
		return
	}

	err := hs.service.VerifyEmail(verificationCode)

	if err != nil {
		log.Printf("Error verifying email with code %s: %v", verificationCode, err)
		switch {
		case errors.Is(err, service.ErrVerificationNotFound):
			http.Error(w, "Verification code not found", http.StatusNotFound) // 404
		case errors.Is(err, service.ErrVerificationExpired):
			http.Error(w, "Verification code has expired", http.StatusBadRequest) // 400
		case errors.Is(err, service.ErrVerificationInvalid):
			http.Error(w, "Verification code is invalid or already used", http.StatusBadRequest) // 400
		case errors.Is(err, service.ErrReferrerNotFound):
			// This is an internal issue, likely shouldn't happen if request flow is right
			http.Error(w, "Associated referrer account not found", http.StatusInternalServerError) // 500
		default:
			// Includes database errors from the service layer
			http.Error(w, "Failed to verify email", http.StatusInternalServerError) // 500
		}
		return
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully."})
}
