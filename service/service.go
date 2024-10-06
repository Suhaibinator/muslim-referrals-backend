package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"muslim-referrals-backend/database"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/oauth2"
)

type Service struct {
	oauthConfig   *oauth2.Config
	userToIdCache *ttlcache.Cache[string, uint64]
	dbDriver      *database.DbDriver
}

func NewService(oauthConfig *oauth2.Config, dbDriver *database.DbDriver) *Service {
	userToIdCache := ttlcache.New[string, uint64](
		ttlcache.WithTTL[string, uint64](24 * time.Hour),
	)
	return &Service{
		oauthConfig:   oauthConfig,
		userToIdCache: userToIdCache,
		dbDriver:      dbDriver,
	}
}

func (s *Service) GetTokenFromCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(ctx, code)
}

func (s *Service) GetUserIdFromTokenDigest(ctx context.Context, tokenDigest string) (uint64, bool, error) {
	result := s.userToIdCache.Get(tokenDigest)
	if result != nil {
		return result.Value(), false, nil
	}

	userInfo, err := s.queryGoogleForEmail(ctx, tokenDigest)
	if err != nil {
		return 0, true, err
	}
	newUser := false
	user := s.dbDriver.GetUserByEmail(userInfo.Email)
	if user == nil {
		newUser = true
		user, err = s.HandleNewUser(ctx, tokenDigest, userInfo)
		if err != nil {
			return 0, newUser, err
		}
	}

	s.userToIdCache.Set(tokenDigest, user.Id, ttlcache.DefaultTTL)

	return user.Id, newUser, nil
}

func (s *Service) queryGoogleForEmail(ctx context.Context, tokenDigest string) (*GoogleUserInfo, error) {

	// Hex decode the token digest into bytes
	tokenBytes, err := hex.DecodeString(tokenDigest)
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
