package service

import (
	"context"

	"github.com/Suhaibinator/muslim-referrals-backend/database"

	"github.com/jellydator/ttlcache/v3"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func (s *Service) HandleNewUser(ctx context.Context, tokenDigest string, userInfo *GoogleUserInfo) (*database.User, error) {
	user := database.User{
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Email:     userInfo.Email,
	}

	createdUser, err := s.dbDriver.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	s.userToIdCache.Set(tokenDigest, createdUser.Id, ttlcache.DefaultTTL)
	return createdUser, nil
}
