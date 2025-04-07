package service

import (
	"errors"
	"fmt"
	"log"
	"muslim-referrals-backend/database"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2" // Added Resend import
	"gorm.io/gorm"
)

var (
	ErrVerificationNotFound     = errors.New("verification code not found")
	ErrVerificationExpired      = errors.New("verification code expired")
	ErrVerificationInvalid      = errors.New("verification status invalid")
	ErrReferrerNotFound         = errors.New("referrer not found for user")
	ErrActiveVerificationExists = errors.New("an active verification request already exists for this email")
	ErrMaxVerificationsReached  = errors.New("maximum number of active verification requests reached for this user")
	ErrEmailSendingDisabled     = errors.New("email sending is currently disabled (missing API key)")
	ErrEmailSendFailed          = errors.New("failed to send verification email")
)

const (
	emailVerificationTTL          = 24 * time.Hour // Verification link valid for 24 hours
	maxActiveVerificationsPerUser = 3              // Max pending requests per user
	// IMPORTANT: Replace with your actual verified sender email and base URL
	verificationSenderEmail = "verify@yourdomain.com" // TODO: Replace with verified sender
	verificationBaseURL     = "http://localhost:8080" // TODO: Replace with actual frontend/API base URL
)

// --- Helper Functions for RequestEmailVerification ---

// checkVerificationPreconditions performs rate limiting and existence checks.
func (s *Service) checkVerificationPreconditions(userID uint64, emailToVerify string) error {
	// 1. Check if an active verification already exists for this email
	existingVerification, err := s.dbDriver.GetActiveVerificationByEmail(emailToVerify)
	if err != nil {
		log.Printf("Error checking for existing verification for email %s: %v", emailToVerify, err)
		return fmt.Errorf("database error checking existing verification: %w", err)
	}
	if existingVerification != nil {
		log.Printf("User %d attempted to verify email %s, but an active request already exists (ID: %s)", userID, emailToVerify, existingVerification.ID)
		return ErrActiveVerificationExists
	}

	// 2. Check user's active verification count (Rate Limit)
	activeCount, err := s.dbDriver.CountActiveVerificationsForUser(userID)
	if err != nil {
		log.Printf("Error counting active verifications for user %d: %v", userID, err)
		return fmt.Errorf("database error checking verification count: %w", err)
	}
	if activeCount >= maxActiveVerificationsPerUser {
		log.Printf("User %d attempted to verify email %s, but reached max active limit (%d)", userID, emailToVerify, maxActiveVerificationsPerUser)
		return ErrMaxVerificationsReached
	}
	return nil
}

// createVerificationRecord creates the initial DB entry for the verification request.
func (s *Service) createVerificationRecord(userID uint64, emailToVerify string) (*database.EmailVerification, error) {
	verificationCode := uuid.NewString()
	expiresAt := time.Now().Add(emailVerificationTTL)
	verification := &database.EmailVerification{
		ID:               verificationCode, // Use UUID as primary key for easier lookup
		Email:            emailToVerify,
		UserID:           userID,
		VerificationCode: verificationCode, // Store code redundantly for potential flexibility
		ExpiresAt:        expiresAt,
		Status:           database.EmailVerificationStatusClaimed, // Initial status
	}

	createdVerification, err := s.dbDriver.CreateEmailVerification(verification)
	if err != nil {
		log.Printf("Error creating email verification record for user %d, email %s: %v", userID, emailToVerify, err)
		return nil, fmt.Errorf("failed to create verification record: %w", err)
	}
	log.Printf("Created email verification record ID %s for user %d, email %s", createdVerification.ID, userID, emailToVerify)
	return createdVerification, nil
}

// sendVerificationEmail handles the construction and sending of the email via Resend.
func (s *Service) sendVerificationEmail(verification *database.EmailVerification) error {
	// Check if API key is configured
	if os.Getenv("RESEND_API_KEY") == "" {
		log.Printf("WARN: Resend API key not set. Skipping email send for verification %s.", verification.ID)
		return ErrEmailSendingDisabled
	}

	verificationLink := fmt.Sprintf("%s/api/email-verification/verify/%s", verificationBaseURL, verification.VerificationCode)
	subject := "Verify Your Email Address"
	htmlBody := fmt.Sprintf(`
		<h1>Welcome to Muslim Referrals!</h1>
		<p>Please verify your email address by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>This link will expire in %s.</p>
		<p>If you did not request this verification, please ignore this email.</p>
	`, verificationLink, emailVerificationTTL.String())

	params := &resend.SendEmailRequest{
		From:    verificationSenderEmail, // Use configured sender
		To:      []string{verification.Email},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := s.resendClient.Emails.Send(params)
	if err != nil {
		log.Printf("ERROR sending verification email via Resend for verification ID %s (User: %d, Email: %s): %v",
			verification.ID, verification.UserID, verification.Email, err)
		return ErrEmailSendFailed // Return generic send error
	}

	log.Printf("Successfully sent verification email via Resend for verification ID %s (User: %d, Email: %s). Resend ID: %s",
		verification.ID, verification.UserID, verification.Email, sent.Id)
	return nil // Send successful
}

// updateVerificationStatusAfterSend updates the DB status based on the email send outcome.
func (s *Service) updateVerificationStatusAfterSend(verification *database.EmailVerification, sendErr error) {
	if verification == nil {
		return // Should not happen, but defensive check
	}

	initialStatus := verification.Status // Store initial status before potential update

	if sendErr != nil {
		verification.Status = database.EmailVerificationStatusSendFailed
	} else {
		verification.Status = database.EmailVerificationStatusSent
	}

	// Avoid unnecessary DB update if status didn't change (e.g., if it was already SendFailed)
	if verification.Status == initialStatus {
		return
	}

	updateErr := s.dbDriver.UpdateEmailVerification(verification)
	if updateErr != nil {
		// Log error based on the intended status update
		if sendErr != nil {
			log.Printf("ERROR updating verification status to SendFailed for ID %s after send error: %v", verification.ID, updateErr)
		} else {
			log.Printf("ERROR updating verification status to Sent for ID %s after successful send: %v", verification.ID, updateErr)
		}
		// Note: The primary operation (sending or failing to send) already determined the outcome.
		// This DB update failure is secondary but should be logged prominently.
	} else {
		log.Printf("Successfully updated verification status for ID %s to %d", verification.ID, verification.Status)
	}
}

// --- Main Service Methods ---

// RequestEmailVerification creates a new email verification request and sends the email.
func (s *Service) RequestEmailVerification(userID uint64, emailToVerify string) error {
	// 1. Perform Pre-checks (Rate limits, existing requests)
	if err := s.checkVerificationPreconditions(userID, emailToVerify); err != nil {
		return err // Return specific errors like ErrActiveVerificationExists, ErrMaxVerificationsReached
	}

	// 2. Create the initial verification record in the database
	verification, err := s.createVerificationRecord(userID, emailToVerify)
	if err != nil {
		return err // Return error from DB creation
	}

	// 3. Attempt to send the verification email
	sendErr := s.sendVerificationEmail(verification)

	// 4. Update the verification status based on send result (handles DB update internally)
	s.updateVerificationStatusAfterSend(verification, sendErr)

	// Return the original send error, if any, to the caller (API handler)
	// This ensures the API layer knows if the primary action (sending) failed.
	return sendErr
}

// VerifyEmail verifies an email address using the provided verification code.
// (No major refactoring needed here for now, structure is reasonably clear)
func (s *Service) VerifyEmail(verificationCode string) error {
	verification, err := s.dbDriver.GetEmailVerificationByCode(verificationCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Verification code not found: %s", verificationCode)
			return ErrVerificationNotFound
		}
		log.Printf("Error retrieving verification code %s: %v", verificationCode, err)
		return fmt.Errorf("database error retrieving verification: %w", err)
	}

	// Check status - Should be 'Sent' to be verifiable
	if verification.Status != database.EmailVerificationStatusSent {
		log.Printf("Verification code %s has invalid status for verification: %d (expected %d)",
			verificationCode, verification.Status, database.EmailVerificationStatusSent)
		// Covers Claimed, Verified, Expired, SendFailed statuses
		return ErrVerificationInvalid
	}

	// Check expiry
	if time.Now().After(verification.ExpiresAt) {
		log.Printf("Verification code %s expired at %s", verificationCode, verification.ExpiresAt)
		verification.Status = database.EmailVerificationStatusExpired
		updateErr := s.dbDriver.UpdateEmailVerification(verification)
		if updateErr != nil {
			log.Printf("Error updating expired verification status for code %s: %v", verificationCode, updateErr)
			// Log error but still return expired error to user
		}
		return ErrVerificationExpired
	}

	// --- Verification Successful ---

	// 1. Update Verification Status
	verification.Status = database.EmailVerificationStatusVerified
	err = s.dbDriver.UpdateEmailVerification(verification)
	if err != nil {
		log.Printf("Error updating verification status to verified for code %s: %v", verificationCode, err)
		return fmt.Errorf("database error updating verification status: %w", err)
	}

	// 2. Update Referrer's Corporate Email
	// GetReferrerByUserId doesn't return an error, it returns nil if not found
	referrer := s.dbDriver.GetReferrerByUserId(verification.UserID)
	if referrer == nil {
		log.Printf("Referrer not found for user ID %d during email verification %s", verification.UserID, verificationCode)
		// This shouldn't happen if the request flow is correct, but handle defensively.
		return ErrReferrerNotFound
	}

	referrer.CorporateEmail = verification.Email
	// UpdateReferrer requires UserID and returns (*Referrer, error)
	_, err = s.dbDriver.UpdateReferrer(verification.UserID, referrer)
	if err != nil {
		log.Printf("Error updating referrer corporate email for user ID %d after verification %s: %v", verification.UserID, verificationCode, err)
		// Critical: Log this, but the verification itself succeeded.
		// Consider how to handle this - maybe retry? For now, return success but log prominently.
		// Or potentially revert the verification status update? Less ideal.
		return fmt.Errorf("database error updating referrer email: %w", err) // Or return nil and just log?
	}

	log.Printf("Successfully verified email %s for user %d (referrer %d) using code %s", verification.Email, verification.UserID, referrer.ReferrerId, verificationCode)
	return nil
}
