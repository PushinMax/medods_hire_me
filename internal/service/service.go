package service

import (
	"medods_hire_me/internal/blacklist"
	"medods_hire_me/internal/mailer"
	"medods_hire_me/internal/repository"
	"medods_hire_me/internal/utils"

	"time"
)

type AuthApi interface {
	IssueToken(userID, ip string) (*utils.TokenPairResponse, error)
	RefreshToken(refreshToken, accessToken, ip string) (*utils.TokenPairResponse, error)
	JWTAuth(token string) (*utils.CustomClaims, error)	
	SetEmail(userID, email string) (error)
}

type Service struct {
	AuthApi
}

func New(repos *repository.Repository, mailer *mailer.Mailer, blacklist *blacklist.Blacklist) *Service {
	return &Service{
		AuthApi: newApiService(repos, mailer, blacklist),
	}
}

type JWTConfig struct {
	AccessSecret  string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}