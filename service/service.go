package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"muslim-referrals-backend/database"
	"os"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/resend/resend-go/v2"
	"golang.org/x/oauth2"
)

type Service struct {
	oauthConfig   *oauth2.Config
	userToIdCache *ttlcache.Cache[string, uint64]
	dbDriver      *database.DbDriver
	resendClient  *resend.Client // Added Resend client
}

func NewService(oauthConfig *oauth2.Config, dbDriver *database.DbDriver) *Service {
	userToIdCache := ttlcache.New[string, uint64](
		ttlcache.WithTTL[string, uint64](24 * time.Hour),
	)

	// Initialize Resend client
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("WARN: RESEND_API_KEY environment variable not set. Email sending will be disabled.")
		// Allow service to start without API key for environments where email isn't needed/configured
	}
	resendClient := resend.NewClient(apiKey) // Client is usable even if apiKey is "" (calls will fail)

	return &Service{
		oauthConfig:   oauthConfig,
		userToIdCache: userToIdCache,
		dbDriver:      dbDriver,
		resendClient:  resendClient, // Store the client
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
