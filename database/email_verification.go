package database

import (
	"time"

	"gorm.io/gorm"
)

// Import gorm explicitly if not already present or if auto-formatter didn't add it globally
// import "gorm.io/gorm" // This line might be redundant depending on formatter/existing imports

type EmailVerificationStatus int

const (
	EmailVerificationStatusClaimed    EmailVerificationStatus = iota // Initial state upon request
	EmailVerificationStatusSent                                      // Email successfully sent
	EmailVerificationStatusVerified                                  // User clicked link successfully
	EmailVerificationStatusExpired                                   // Link expired before verification
	EmailVerificationStatusSendFailed                                // Attempt to send email failed
)

type EmailVerification struct {
	ID               string                  `json:"id" gorm:"primaryKey"`
	Email            string                  `json:"email" gorm:"unique;not null"`
	UserID           uint64                  `json:"user_id" gorm:"not null;index"`   // Added UserID field
	User             User                    `gorm:"foreignKey:UserID;references:Id"` // Added User relationship
	VerificationCode string                  `json:"verification_code" gorm:"not null"`
	ExpiresAt        time.Time               `json:"expires_at" gorm:"not null"`
	Status           EmailVerificationStatus `json:"status" gorm:"not null;default:0"` // Default to Claimed
}

func (db *DbDriver) CreateEmailVerification(record *EmailVerification) (*EmailVerification, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) GetEmailVerificationByCode(code string) (*EmailVerification, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var verification EmailVerification
	result := db.db.Where("verification_code = ?", code).First(&verification)
	if result.Error != nil {
		return nil, result.Error // Could be gorm.ErrRecordNotFound or other error
	}
	return &verification, nil
}

func (db *DbDriver) UpdateEmailVerification(record *EmailVerification) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Use Save to update all fields, including Status and potentially ExpiresAt if needed later
	result := db.db.Save(record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetActiveVerificationByEmail finds an existing, unexpired verification request
// for a specific email that is still in a pending state (Claimed or Sent).
// Returns the verification record if found, nil otherwise. Error indicates a DB issue.
func (db *DbDriver) GetActiveVerificationByEmail(email string) (*EmailVerification, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var verification EmailVerification
	now := time.Now()
	// Check for requests that are not yet verified/expired/failed and haven't passed expiry time
	result := db.db.Where("email = ? AND status IN (?, ?) AND expires_at > ?",
		email, EmailVerificationStatusClaimed, EmailVerificationStatusSent, now).First(&verification)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not a DB error in this context
		}
		return nil, result.Error // Actual database error
	}
	return &verification, nil
}

// CountActiveVerificationsForUser counts the number of unexpired verification requests
// for a specific user that are still in a pending state (Claimed or Sent).
func (db *DbDriver) CountActiveVerificationsForUser(userID uint64) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var count int64
	now := time.Now()
	// Count requests for the user that are not yet verified/expired/failed and haven't passed expiry time
	result := db.db.Model(&EmailVerification{}).Where("user_id = ? AND status IN (?, ?) AND expires_at > ?",
		userID, EmailVerificationStatusClaimed, EmailVerificationStatusSent, now).Count(&count)

	if result.Error != nil {
		return 0, result.Error // Database error
	}
	return count, nil
}
