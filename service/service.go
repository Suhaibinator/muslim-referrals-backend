package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/Suhaibinator/muslim-referrals-backend/database"

	"github.com/jellydator/ttlcache/v3"
	"github.com/resend/resend-go/v2"
	"golang.org/x/oauth2"
)

// DatabaseOperations defines the interface for database interactions needed by the service.
// This allows mocking the database layer for unit tests.
type DatabaseOperations interface {
	// Email Verification Methods
	GetActiveVerificationByEmail(email string) (*database.EmailVerification, error)
	CountActiveVerificationsForUser(userID uint64) (int64, error)
	CreateEmailVerification(verification *database.EmailVerification) (*database.EmailVerification, error)
	UpdateEmailVerification(verification *database.EmailVerification) error
	GetEmailVerificationByCode(code string) (*database.EmailVerification, error)

	// Referrer Methods
	GetReferrerByUserId(userID uint64) *database.Referrer
	UpdateReferrer(userID uint64, referrer *database.Referrer) (*database.Referrer, error)

	// User Methods
	GetUserByEmail(email string) *database.User
	CreateUser(user *database.User) (*database.User, error)
	// Add other DB methods used by the service here...
}

// EmailSender defines the interface for sending emails, matching resend's EmailsSvc.
// This allows mocking the email sending functionality.
type EmailSender interface {
	Send(params *resend.SendEmailRequest) (*resend.SendEmailResponse, error)
	// Add other resend.EmailsSvc methods if needed by the service
}

type Service struct {
	oauthConfig   *oauth2.Config
	userToIdCache *ttlcache.Cache[string, uint64]
	dbDriver      DatabaseOperations // Use the interface type
	emailSender   EmailSender        // Use the interface type (can be resend.EmailsSvc)
}

// SetUserIDForToken allows tests to seed the cache with a token to user ID mapping.
func (s *Service) SetUserIDForToken(token string, userID uint64) {
	s.userToIdCache.Set(token, userID, ttlcache.DefaultTTL)
}

// NewService now accepts interfaces for dependencies, improving testability.
func NewService(oauthConfig *oauth2.Config, dbDriver DatabaseOperations, emailSender EmailSender) *Service {
	userToIdCache := ttlcache.New[string, uint64](
		ttlcache.WithTTL[string, uint64](24 * time.Hour),
	)

	// Dependencies (dbDriver, emailSender) are now injected.
	// No need to initialize Resend client here; it's passed in.

	return &Service{
		oauthConfig:   oauthConfig,
		userToIdCache: userToIdCache,
		dbDriver:      dbDriver,    // Assign injected DB interface
		emailSender:   emailSender, // Assign injected email sender interface
	}
}

func (s *Service) GetTokenFromCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(ctx, code)
}

func (s *Service) GetUserIdFromTokenDigest(ctx context.Context, tokenDigest string) (uint64, bool, error) {
	result := s.userToIdCache.Get(tokenDigest)
	if result != nil {
		log.Printf("[GetUserIdFromTokenDigest] Cache hit for token digest %s", tokenDigest)
		return result.Value(), false, nil
	}
	log.Printf("[GetUserIdFromTokenDigest] Cache miss for token digest %s", tokenDigest)

	userInfo, err := s.queryGoogleForEmail(ctx, tokenDigest)
	log.Printf("[GetUserIdFromTokenDigest] Got user info: %+v", userInfo)
	if err != nil {
		log.Printf("[GetUserIdFromTokenDigest] Error getting user info: %v", err)
		return 0, true, err
	}

	newUser := false
	user := s.dbDriver.GetUserByEmail(userInfo.Email)
	if user == nil {
		newUser = true
		user, err = s.HandleNewUser(ctx, tokenDigest, userInfo)
		if err != nil {
			log.Printf("[GetUserIdFromTokenDigest] Error handling new user: %v", err)
			return 0, newUser, err
		}
	}

	s.userToIdCache.Set(tokenDigest, user.Id, ttlcache.DefaultTTL)
	log.Printf("[GetUserIdFromTokenDigest] Set cache for token digest %s", tokenDigest)

	return user.Id, newUser, nil
}

func (s *Service) queryGoogleForEmail(ctx context.Context, tokenDigest string) (*GoogleUserInfo, error) {

	// B64 decode the token digest into bytes
	tokenBytes, err := base64.StdEncoding.DecodeString(tokenDigest)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	unmarshalErr := json.Unmarshal(tokenBytes, &token)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	// Get user information from Google
	client := s.oauthConfig.Client(ctx, &token)
	resp, obtainUserInfoErr := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if obtainUserInfoErr != nil {
		return nil, obtainUserInfoErr
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	decodeErr := json.NewDecoder(resp.Body).Decode(&userInfo)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return &userInfo, nil
}
