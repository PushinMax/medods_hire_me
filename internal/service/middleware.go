package service

import (
	"medods_hire_me/internal/utils"

	"errors"
	"log"
)

func (s *ApiService) JWTAuth(token string) (*utils.CustomClaims, error) {
	claims, err := utils.ValidateJWT(token, s.cfg.AccessSecret)
	if err != nil {
		return nil, err
	}
	if ip, ok := s.blacklist.ContainsAndGetIp(claims.JTI); ok {
		if ip != claims.IP {
			email, err := s.repo.GetEmail(claims.Subject)
			if err != nil {
				log.Printf("User %s did not specify email.", claims.Subject)
			}
			_ = s.mailer.Send(email, "Warning", "A suspicious transaction has been detected")
			
			return nil, errors.New("using token from unknown IP")
		}
		return claims, nil
	}
	return nil, errors.New("token in blacklist")
} 