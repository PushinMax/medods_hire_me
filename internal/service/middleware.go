package service

import (
	"medods_hire_me/internal/utils"
)

func (s *ApiService) JWTAuth(token string) (*utils.CustomClaims, error) {
	return utils.ValidateJWT(token, s.cfg.AccessSecret)
} 