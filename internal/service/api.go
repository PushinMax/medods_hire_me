package service

import (
	"encoding/base64"
	"log"
	"medods_hire_me/internal/mailer"
	"medods_hire_me/internal/repository"
	"medods_hire_me/internal/utils"

	"fmt"

	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"errors"
	"os"

	"github.com/spf13/viper"
)

type ApiService struct {
	repo repository.RepoApi
	cfg  JWTConfig
	mailer mailer.MailApi
}

func newApiService(repo repository.RepoApi, mailer mailer.MailApi) *ApiService {
	accessExpiry, _ := time.ParseDuration(viper.GetString("jwt.access_expiry"))
	refreshExpiry, _ := time.ParseDuration(viper.GetString("jwt.refresh_expiry"))
	return &ApiService{
		repo: repo,
		cfg: JWTConfig{
			AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		mailer: mailer,
	}
}

func (s *ApiService) IssueToken(userID, ip string) (*utils.TokenPairResponse, error) {
	if err := s.repo.Exists(userID); err != nil {
		return nil, fmt.Errorf("user check failed: %w", err)
	}

	jti := uuid.New().String()

	accessToken, err := utils.GenerateJWT(
		userID,
		ip,
		jti,
		s.cfg.AccessSecret,
		s.cfg.AccessExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("access token generation failed: %w", err)
	}

	rawRefreshToken, err := utils.GenerateRandomToken(64)
	if err != nil {
		return nil, err
	}
	

	refreshToken := base64.StdEncoding.EncodeToString(rawRefreshToken)

	hashedRefresh, err := bcrypt.GenerateFromPassword(rawRefreshToken, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("refresh token hash failed: %w", err)
	}

	expiryTime := time.Now().Add(s.cfg.RefreshExpiry)
	if err := s.repo.SaveRefreshToken(
		userID,
		string(hashedRefresh),
		ip,
		jti,
		expiryTime,
	); err != nil {
		return nil, fmt.Errorf("token save failed: %w", err)
	}

	return &utils.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ApiService) RefreshToken(refreshToken, accessToken, ip string) (*utils.TokenPairResponse, error) {
	rawToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	rawAccessToken, err := utils.ValidateJWT(accessToken, s.cfg.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("bad format of accesToken: %s", err.Error())
	}

	if err := s.repo.TokenVerification(&repository.TokenData{
		RefreshToken: string(rawToken),
		AccessToken:  rawAccessToken,
		IP:           ip,
	}); err != nil {
		if err.Error() == "IP do not match" {
			email, err := s.repo.GetEmail(rawAccessToken.Subject)
			if err != nil {
				log.Printf("User %s did not specify email.", rawAccessToken.Subject)
			}
			s.mailer.Send(email, "Warning", "A suspicious transaction has been detected. The above token was accessed by IP:" + rawAccessToken.IP)
			return nil, errors.New("the token is not provided for your IP")
		}
		return nil, fmt.Errorf("TokenVerification: %s", err.Error())
	}
	return s.IssueToken(rawAccessToken.Subject, ip)
}

func (s *ApiService) SetEmail(userID, email string) error {
	if oldEmail, err:= s.repo.GetEmail(userID); err == nil {
		s.mailer.Send(oldEmail, "Last notification", "The notification email has been changed.")
	} 
	return s.repo.SetEmail(userID, email)
}