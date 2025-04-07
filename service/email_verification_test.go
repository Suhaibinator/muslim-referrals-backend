package service_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Suhaibinator/muslim-referrals-backend/database"
	"github.com/Suhaibinator/muslim-referrals-backend/service"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// --- Mock Database Driver ---

// MockDatabaseDriver needs to implement the interface expected by the service.
// Let's assume there's an interface like `database.DbOperations` that it needs to satisfy.
// For now, we'll just embed mock.Mock.
type MockDatabaseDriver struct {
	mock.Mock
}

// Ensure MockDatabaseDriver implements the necessary methods (adjust interface name if needed)
var _ service.DatabaseOperations = (*MockDatabaseDriver)(nil) // Check interface implementation

func (m *MockDatabaseDriver) GetActiveVerificationByEmail(email string) (*database.EmailVerification, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Directly return the configured value, asserting its type.
	retVal, ok := args.Get(0).(*database.EmailVerification)
	if !ok && args.Get(0) != nil { // Allow nil return without panic
		// Panic if the test setup provided the wrong type in Return()
		panic(fmt.Sprintf("mock GetActiveVerificationByEmail: Return(0) is not *database.EmailVerification: %T", args.Get(0)))
	}
	return retVal, args.Error(1)
}

// Restore the correct implementation for CountActiveVerificationsForUser
func (m *MockDatabaseDriver) CountActiveVerificationsForUser(userID uint64) (int64, error) {
	args := m.Called(userID)
	// Directly return the configured value, asserting its type.
	retVal, ok := args.Get(0).(int64)
	if !ok {
		// Panic if the test setup provided the wrong type in Return()
		panic(fmt.Sprintf("mock CountActiveVerificationsForUser: Return(0) is not int64: %T", args.Get(0)))
	}
	return retVal, args.Error(1)
}

func (m *MockDatabaseDriver) CreateEmailVerification(verification *database.EmailVerification) (*database.EmailVerification, error) {
	args := m.Called(verification)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Directly return the configured value, asserting its type.
	retVal, ok := args.Get(0).(*database.EmailVerification)
	if !ok {
		// Panic if the test setup provided the wrong type in Return()
		panic(fmt.Sprintf("mock CreateEmailVerification: Return(0) is not *database.EmailVerification: %T", args.Get(0)))
	}
	return retVal, args.Error(1)
}

func (m *MockDatabaseDriver) UpdateEmailVerification(verification *database.EmailVerification) error {
	args := m.Called(verification)
	return args.Error(0)
}

func (m *MockDatabaseDriver) GetEmailVerificationByCode(code string) (*database.EmailVerification, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.EmailVerification), args.Error(1)
}

func (m *MockDatabaseDriver) GetReferrerByUserId(userID uint64) *database.Referrer {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*database.Referrer)
}

func (m *MockDatabaseDriver) UpdateReferrer(userID uint64, referrer *database.Referrer) (*database.Referrer, error) {
	args := m.Called(userID, referrer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.Referrer), args.Error(1)
}

// Add missing methods required by service.DatabaseOperations
func (m *MockDatabaseDriver) GetUserByEmail(email string) *database.User {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*database.User)
}

func (m *MockDatabaseDriver) CreateUser(user *database.User) (*database.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Directly return the configured value, asserting its type.
	retVal, ok := args.Get(0).(*database.User)
	if !ok {
		// Panic if the test setup provided the wrong type in Return()
		panic(fmt.Sprintf("mock CreateUser: Return(0) is not *database.User: %T", args.Get(0)))
	}
	return retVal, args.Error(1)
}

// --- Mock Resend Client ---

// Interface for the part of Resend client we use
type ResendEmailsAPI interface {
	Send(params *resend.SendEmailRequest) (*resend.SendEmailResponse, error)
}

type MockResendEmailsAPI struct {
	mock.Mock
}

// Note: We remove the interface satisfaction check `var _ resend.EmailsSvc = (*MockResendEmailsAPI)(nil)`
// because fully implementing the interface is causing issues with potentially missing/incorrectly named
// methods like `Update`. Since the code under test only uses `Send`, we only mock that.

func (m *MockResendEmailsAPI) Send(params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.SendEmailResponse), args.Error(1)
}

// We only mock SendWithContext if it's actually used by the service logic.
// Assuming email_verification.go uses the non-context Send method based on previous reads.
// func (m *MockResendEmailsAPI) SendWithContext(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
// 	return m.Send(params)
// }

// --- Test Setup ---

// Helper to create a service instance with mocks.
// Allows optionally passing a specific EmailSender (e.g., nil for testing disabled state).
// If emailSenderOverride is nil, the service instance will receive nil for its emailSender field.
func setupServiceWithMocks(emailSenderOverride service.EmailSender) (*service.Service, *MockDatabaseDriver, *MockResendEmailsAPI) {
	mockDB := new(MockDatabaseDriver)
	// Always create the mockResendEmails instance so the caller can potentially set expectations on it,
	// even if the service itself receives nil.
	mockResendEmails := new(MockResendEmailsAPI)

	// Create a dummy OAuth config (replace with actual config if needed for other tests)
	dummyOauthConfig := &oauth2.Config{
		ClientID:     "dummy-client-id",
		ClientSecret: "dummy-client-secret",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	// Instantiate the real service using the constructor.
	// NewService takes *oauth2.Config and *database.DbDriver.
	// It initializes its own internal resend client.
	// We need to pass mockDB, assuming it satisfies the interface expected by *database.DbDriver methods used by the service.
	// If *database.DbDriver is a struct, mockDB needs to implement an interface that DbDriver implicitly satisfies,
	// or we need a more complex setup.
	// Instantiate the service using the refactored constructor, injecting mocks.
	// MockDatabaseDriver implicitly satisfies the DatabaseOperations interface (as long as methods match).
	// MockResendEmailsAPI implicitly satisfies the EmailSender interface.
	// Pass the override directly. If it's nil, the service gets nil.
	s := service.NewService(dummyOauthConfig, mockDB, emailSenderOverride)

	// Now the service `s` has the mocks injected correctly (or nil sender).
	// Return the service, the db mock, and the *original* email mock instance
	// so the caller can set expectations on it if needed.
	// Ensure RESEND_API_KEY is set for tests that might still check the env var directly
	// (though the refactored service shouldn't rely on it internally anymore for the client instance).
	// in tests where sending *should* be attempted.
	os.Setenv("RESEND_API_KEY", "test-key-for-testing")

	return s, mockDB, mockResendEmails
}

// --- Test Cases for RequestEmailVerification ---

func TestRequestEmailVerification_Success(t *testing.T) {
	// Create the mock sender first
	mockResendEmails := new(MockResendEmailsAPI)
	// Pass the created mock to the setup function
	// Note: We ignore the third return value from setupServiceWithMocks as we are using the one we created.
	s, mockDB, _ := setupServiceWithMocks(mockResendEmails)
	userID := uint64(1)
	email := "test@example.com"
	testUUID := uuid.NewString()

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), nil).Once()
	// Expect CreateEmailVerification and capture the argument to set the ID
	mockDB.On("CreateEmailVerification", mock.AnythingOfType("*database.EmailVerification")).Run(func(args mock.Arguments) {
		ver := args.Get(0).(*database.EmailVerification)
		ver.ID = testUUID                                    // Ensure the mock returns the expected ID for later checks
		ver.Status = database.EmailVerificationStatusClaimed // Set initial status
		assert.Equal(t, email, ver.Email)
		assert.Equal(t, userID, ver.UserID)
		assert.WithinDuration(t, time.Now().Add(24*time.Hour), ver.ExpiresAt, 5*time.Second) // Check expiry
	}).Return(&database.EmailVerification{ID: testUUID, Email: email, UserID: userID, Status: database.EmailVerificationStatusClaimed /* Add other fields if needed by subsequent code */}, nil).Once() // Return a concrete object
	mockResendEmails.On("Send", mock.AnythingOfType("*resend.SendEmailRequest")).Return(&resend.SendEmailResponse{Id: "resend-id"}, nil).Once()
	// Expect UpdateEmailVerification with status Sent
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(ver *database.EmailVerification) bool {
		return ver.ID == testUUID && ver.Status == database.EmailVerificationStatusSent
	})).Return(nil).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockResendEmails.AssertExpectations(t)
}

func TestRequestEmailVerification_ErrActiveVerificationExists(t *testing.T) {
	// Pass nil override, email sender not used, use _ for mockResendEmails
	s, mockDB, _ := setupServiceWithMocks(nil)
	userID := uint64(1)
	email := "test@example.com"
	existingVerification := &database.EmailVerification{ID: "existing-uuid", Email: email, UserID: userID, Status: database.EmailVerificationStatusSent} // Example existing

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(existingVerification, nil).Once()
	// No other mocks should be called

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrActiveVerificationExists, err)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "CountActiveVerificationsForUser", mock.Anything)
	mockDB.AssertNotCalled(t, "CreateEmailVerification", mock.Anything)
	// No need for mockResendEmails assertions here
	mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
}

func TestRequestEmailVerification_ErrMaxVerificationsReached(t *testing.T) {
	// Pass nil override, email sender not used, use _
	s, mockDB, _ := setupServiceWithMocks(nil)
	userID := uint64(1)
	email := "test@example.com"

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(3), nil).Once() // Assuming max is 3
	// No other mocks should be called

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrMaxVerificationsReached, err)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "CreateEmailVerification", mock.Anything)
	// No need for mockResendEmails assertions here
	mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
}

func TestRequestEmailVerification_CheckPreconditionsDbError_GetActive(t *testing.T) {
	// Pass nil override, email sender not used, use _
	s, mockDB, _ := setupServiceWithMocks(nil)
	userID := uint64(1)
	email := "test@example.com"
	dbError := errors.New("db get active error")

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, dbError).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error checking existing verification")
	assert.ErrorIs(t, err, dbError)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "CountActiveVerificationsForUser", mock.Anything)
	mockDB.AssertNotCalled(t, "CreateEmailVerification", mock.Anything)
	// No need for mockResendEmails assertions here
	mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
}

func TestRequestEmailVerification_CheckPreconditionsDbError_CountActive(t *testing.T) {
	// Pass nil override, email sender not used, use _
	s, mockDB, _ := setupServiceWithMocks(nil)
	userID := uint64(1)
	email := "test@example.com"
	dbError := errors.New("db count active error")

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), dbError).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error checking verification count")
	assert.ErrorIs(t, err, dbError)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "CreateEmailVerification", mock.Anything)
	// No need for mockResendEmails assertions here
	mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
}

func TestRequestEmailVerification_CreateRecordFails(t *testing.T) {
	// Pass nil override, email sender not used, use _
	s, mockDB, _ := setupServiceWithMocks(nil)
	userID := uint64(1)
	email := "test@example.com"
	dbError := errors.New("db create error")

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), nil).Once()
	mockDB.On("CreateEmailVerification", mock.AnythingOfType("*database.EmailVerification")).Return(nil, dbError).Once()
	// No email send or update should happen

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create verification record") // Check if wrapped error contains expected text
	assert.ErrorIs(t, err, dbError)                                         // Check if it wraps the specific DB error
	mockDB.AssertExpectations(t)
	// No need for mockResendEmails assertions here
	mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
}

func TestRequestEmailVerification_EmailSendingDisabled(t *testing.T) {
	// Temporarily unset API key
	originalKey, keyExists := os.LookupEnv("RESEND_API_KEY")
	os.Unsetenv("RESEND_API_KEY")
	defer func() {
		if keyExists {
			os.Setenv("RESEND_API_KEY", originalKey)
		} else {
			os.Unsetenv("RESEND_API_KEY")
		}
	}() // Restore original state

	// Setup service with a nil email sender for this specific test by passing nil override
	s, mockDB, _ := setupServiceWithMocks(nil)
	userID := uint64(1)
	email := "test@example.com"
	testUUID := uuid.NewString()

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), nil).Once()
	// Expect CreateEmailVerification and capture the argument
	mockDB.On("CreateEmailVerification", mock.AnythingOfType("*database.EmailVerification")).Run(func(args mock.Arguments) {
		ver := args.Get(0).(*database.EmailVerification)
		ver.ID = testUUID
		ver.Status = database.EmailVerificationStatusClaimed
	}).Return(&database.EmailVerification{ID: testUUID, Email: email, UserID: userID, Status: database.EmailVerificationStatusClaimed}, nil).Once() // Return a concrete object
	// Email send should NOT be called
	// Expect UpdateEmailVerification with status SendFailed
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(ver *database.EmailVerification) bool {
		return ver.ID == testUUID && ver.Status == database.EmailVerificationStatusSendFailed
	})).Return(nil).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrEmailSendingDisabled, err) // Expect specific error
	mockDB.AssertExpectations(t)
	// No need to assert mockResendEmails was not called, as the service instance has a nil sender.
	// If the code incorrectly tries to call Send on the nil sender, it will panic.
}

func TestRequestEmailVerification_EmailSendFails(t *testing.T) {
	// Create mock sender first
	mockResendEmails := new(MockResendEmailsAPI)
	// Pass the created mock to the setup function
	s, mockDB, _ := setupServiceWithMocks(mockResendEmails)
	userID := uint64(1)
	email := "test@example.com"
	testUUID := uuid.NewString()
	sendError := errors.New("resend failed")

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), nil).Once()
	// Expect CreateEmailVerification and capture the argument
	mockDB.On("CreateEmailVerification", mock.AnythingOfType("*database.EmailVerification")).Run(func(args mock.Arguments) {
		ver := args.Get(0).(*database.EmailVerification)
		ver.ID = testUUID
		ver.Status = database.EmailVerificationStatusClaimed
	}).Return(&database.EmailVerification{ID: testUUID, Email: email, UserID: userID, Status: database.EmailVerificationStatusClaimed}, nil).Once() // Return a concrete object
	mockResendEmails.On("Send", mock.AnythingOfType("*resend.SendEmailRequest")).Return(nil, sendError).Once()
	// Expect UpdateEmailVerification with status SendFailed
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(ver *database.EmailVerification) bool {
		return ver.ID == testUUID && ver.Status == database.EmailVerificationStatusSendFailed
	})).Return(nil).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrEmailSendFailed, err) // Expect specific error
	mockDB.AssertExpectations(t)
	mockResendEmails.AssertExpectations(t)
}

func TestRequestEmailVerification_UpdateStatusAfterSendFails(t *testing.T) {
	// Create mock sender first
	mockResendEmails := new(MockResendEmailsAPI)
	// Pass the created mock to the setup function
	s, mockDB, _ := setupServiceWithMocks(mockResendEmails)
	userID := uint64(1)
	email := "test@example.com"
	testUUID := uuid.NewString()
	updateError := errors.New("db update failed")

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), nil).Once()
	mockDB.On("CreateEmailVerification", mock.AnythingOfType("*database.EmailVerification")).Run(func(args mock.Arguments) {
		ver := args.Get(0).(*database.EmailVerification)
		ver.ID = testUUID
		ver.Status = database.EmailVerificationStatusClaimed
	}).Return(&database.EmailVerification{ID: testUUID, Email: email, UserID: userID, Status: database.EmailVerificationStatusClaimed}, nil).Once() // Return a concrete object
	mockResendEmails.On("Send", mock.AnythingOfType("*resend.SendEmailRequest")).Return(&resend.SendEmailResponse{Id: "resend-id"}, nil).Once()
	// Expect UpdateEmailVerification to fail
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(ver *database.EmailVerification) bool {
		return ver.ID == testUUID && ver.Status == database.EmailVerificationStatusSent
	})).Return(updateError).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	// The primary operation (sending email) succeeded, so RequestEmailVerification should return nil.
	// The failure to update the status is logged internally but doesn't fail the overall request.
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockResendEmails.AssertExpectations(t)
	// We expect logs indicating the update failure, but can't easily assert on logs here.
}

func TestRequestEmailVerification_UpdateStatusAfterSendFail_WhenSendFailed(t *testing.T) {
	// Create mock sender first
	mockResendEmails := new(MockResendEmailsAPI)
	// Pass the created mock to the setup function
	s, mockDB, _ := setupServiceWithMocks(mockResendEmails)
	userID := uint64(1)
	email := "test@example.com"
	testUUID := uuid.NewString()
	sendError := errors.New("resend failed")
	updateError := errors.New("db update failed")

	// Mock expectations
	mockDB.On("GetActiveVerificationByEmail", email).Return(nil, nil).Once()
	mockDB.On("CountActiveVerificationsForUser", userID).Return(int64(0), nil).Once()
	mockDB.On("CreateEmailVerification", mock.AnythingOfType("*database.EmailVerification")).Run(func(args mock.Arguments) {
		ver := args.Get(0).(*database.EmailVerification)
		ver.ID = testUUID
		ver.Status = database.EmailVerificationStatusClaimed
	}).Return(&database.EmailVerification{ID: testUUID, Email: email, UserID: userID, Status: database.EmailVerificationStatusClaimed}, nil).Once() // Return a concrete object
	mockResendEmails.On("Send", mock.AnythingOfType("*resend.SendEmailRequest")).Return(nil, sendError).Once()
	// Expect UpdateEmailVerification to fail (trying to set SendFailed status)
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(ver *database.EmailVerification) bool {
		return ver.ID == testUUID && ver.Status == database.EmailVerificationStatusSendFailed
	})).Return(updateError).Once()

	// Call the function
	err := s.RequestEmailVerification(userID, email)

	// Assertions
	// The primary operation (sending email) failed, so RequestEmailVerification should return ErrEmailSendFailed.
	assert.Error(t, err)
	assert.Equal(t, service.ErrEmailSendFailed, err)
	mockDB.AssertExpectations(t)
	mockResendEmails.AssertExpectations(t)
	// We expect logs indicating the update failure, but can't easily assert on logs here.
}

// --- Test Cases for VerifyEmail ---

func TestVerifyEmail_Success(t *testing.T) {
	// Use default setup (email sender mock instance is created but not needed/used here)
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "valid-code"
	userID := uint64(5)
	email := "verified@example.com"
	referrerID := uint64(10)

	verification := &database.EmailVerification{
		ID:               verificationCode,
		Email:            email,
		UserID:           userID,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(1 * time.Hour),        // Not expired
		Status:           database.EmailVerificationStatusSent, // Correct initial status
	}
	referrer := &database.Referrer{
		ReferrerId:     referrerID,
		UserId:         userID,
		CorporateEmail: "old@example.com", // Initial email
	}

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
	// Expect update to Verified status
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(v *database.EmailVerification) bool {
		return v.ID == verificationCode && v.Status == database.EmailVerificationStatusVerified
	})).Return(nil).Once()
	mockDB.On("GetReferrerByUserId", userID).Return(referrer, nil).Once()
	// Expect update with new email
	mockDB.On("UpdateReferrer", userID, mock.MatchedBy(func(r *database.Referrer) bool {
		return r.ReferrerId == referrerID && r.CorporateEmail == email
	})).Return(referrer, nil).Once()

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestVerifyEmail_NotFound(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "invalid-code"

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(nil, gorm.ErrRecordNotFound).Once()

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrVerificationNotFound, err)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
	mockDB.AssertNotCalled(t, "GetReferrerByUserId", mock.Anything)
	mockDB.AssertNotCalled(t, "UpdateReferrer", mock.Anything, mock.Anything)
}

func TestVerifyEmail_DbErrorOnGet(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "valid-code"
	dbError := errors.New("db get error")

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(nil, dbError).Once()

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error retrieving verification")
	assert.ErrorIs(t, err, dbError)
	mockDB.AssertExpectations(t)
}

func TestVerifyEmail_InvalidStatus(t *testing.T) {
	testCases := []struct {
		name          string
		initialStatus database.EmailVerificationStatus
	}{
		{"Claimed", database.EmailVerificationStatusClaimed},
		{"Verified", database.EmailVerificationStatusVerified},
		{"Expired", database.EmailVerificationStatusExpired},
		{"SendFailed", database.EmailVerificationStatusSendFailed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use default setup
			s, mockDB, _ := setupServiceWithMocks(nil)
			verificationCode := "code-with-bad-status"
			verification := &database.EmailVerification{
				ID:               verificationCode,
				Email:            "test@example.com",
				UserID:           uint64(1),
				VerificationCode: verificationCode,
				ExpiresAt:        time.Now().Add(1 * time.Hour),
				Status:           tc.initialStatus, // Set the invalid status
			}

			// Mock expectations
			mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
			// No update should happen

			// Call the function
			err := s.VerifyEmail(verificationCode)

			// Assertions
			assert.Error(t, err)
			assert.Equal(t, service.ErrVerificationInvalid, err)
			mockDB.AssertExpectations(t)
			mockDB.AssertNotCalled(t, "UpdateEmailVerification", mock.Anything)
			mockDB.AssertNotCalled(t, "GetReferrerByUserId", mock.Anything)
			mockDB.AssertNotCalled(t, "UpdateReferrer", mock.Anything, mock.Anything)
		})
	}
}

func TestVerifyEmail_Expired(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "expired-code"
	verification := &database.EmailVerification{
		ID:               verificationCode,
		Email:            "test@example.com",
		UserID:           uint64(1),
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(-1 * time.Hour),       // Expired
		Status:           database.EmailVerificationStatusSent, // Was sent, but now expired
	}

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
	// Expect update to Expired status
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(v *database.EmailVerification) bool {
		return v.ID == verificationCode && v.Status == database.EmailVerificationStatusExpired
	})).Return(nil).Once() // Assume update succeeds

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrVerificationExpired, err)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetReferrerByUserId", mock.Anything)
	mockDB.AssertNotCalled(t, "UpdateReferrer", mock.Anything, mock.Anything)
}

func TestVerifyEmail_Expired_UpdateFails(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "expired-code"
	updateError := errors.New("db update expired failed")
	verification := &database.EmailVerification{
		ID:               verificationCode,
		Email:            "test@example.com",
		UserID:           uint64(1),
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(-1 * time.Hour), // Expired
		Status:           database.EmailVerificationStatusSent,
	}

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
	// Expect update to Expired status, but it fails
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(v *database.EmailVerification) bool {
		return v.ID == verificationCode && v.Status == database.EmailVerificationStatusExpired
	})).Return(updateError).Once()

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	// Still returns the primary error (Expired) even if the status update fails
	assert.Equal(t, service.ErrVerificationExpired, err)
	mockDB.AssertExpectations(t)
	// Log message expected, but can't assert easily
}

func TestVerifyEmail_UpdateVerificationStatusFails(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "valid-code-update-fail"
	userID := uint64(5)
	email := "verified@example.com"
	updateError := errors.New("db update verified failed")

	verification := &database.EmailVerification{
		ID:               verificationCode,
		Email:            email,
		UserID:           userID,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(1 * time.Hour),
		Status:           database.EmailVerificationStatusSent,
	}

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
	// Expect update to Verified status, but it fails
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(v *database.EmailVerification) bool {
		return v.ID == verificationCode && v.Status == database.EmailVerificationStatusVerified
	})).Return(updateError).Once()
	// No referrer lookups/updates should happen if the first update fails

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error updating verification status")
	assert.ErrorIs(t, err, updateError)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetReferrerByUserId", mock.Anything)
	mockDB.AssertNotCalled(t, "UpdateReferrer", mock.Anything, mock.Anything)
}

func TestVerifyEmail_ReferrerNotFound(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "valid-code-no-referrer"
	userID := uint64(6) // User exists for verification, but not as referrer
	email := "verified@example.com"

	verification := &database.EmailVerification{
		ID:               verificationCode,
		Email:            email,
		UserID:           userID,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(1 * time.Hour),
		Status:           database.EmailVerificationStatusSent,
	}

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(v *database.EmailVerification) bool {
		return v.ID == verificationCode && v.Status == database.EmailVerificationStatusVerified
	})).Return(nil).Once()
	// Referrer lookup returns nil
	mockDB.On("GetReferrerByUserId", userID).Return(nil).Once()
	// No UpdateReferrer call expected

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, service.ErrReferrerNotFound, err)
	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "UpdateReferrer", mock.Anything, mock.Anything)
}

func TestVerifyEmail_UpdateReferrerFails(t *testing.T) {
	// Use default setup
	s, mockDB, _ := setupServiceWithMocks(nil)
	verificationCode := "valid-code-referrer-update-fail"
	userID := uint64(7)
	email := "verified@example.com"
	referrerID := uint64(11)
	updateError := errors.New("db update referrer failed")

	verification := &database.EmailVerification{
		ID:               verificationCode,
		Email:            email,
		UserID:           userID,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().Add(1 * time.Hour),
		Status:           database.EmailVerificationStatusSent,
	}
	referrer := &database.Referrer{
		ReferrerId:     referrerID,
		UserId:         userID,
		CorporateEmail: "old@example.com",
	}

	// Mock expectations
	mockDB.On("GetEmailVerificationByCode", verificationCode).Return(verification, nil).Once()
	mockDB.On("UpdateEmailVerification", mock.MatchedBy(func(v *database.EmailVerification) bool {
		return v.ID == verificationCode && v.Status == database.EmailVerificationStatusVerified
	})).Return(nil).Once()
	mockDB.On("GetReferrerByUserId", userID).Return(referrer, nil).Once()
	// Expect UpdateReferrer, but it fails
	mockDB.On("UpdateReferrer", userID, mock.MatchedBy(func(r *database.Referrer) bool {
		return r.ReferrerId == referrerID && r.CorporateEmail == email
	})).Return(nil, updateError).Once()

	// Call the function
	err := s.VerifyEmail(verificationCode)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error updating referrer email")
	assert.ErrorIs(t, err, updateError)
	mockDB.AssertExpectations(t)
}
