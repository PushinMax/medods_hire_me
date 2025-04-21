package utils

import (
	"github.com/google/uuid"
)

type TokenRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required,uuid"`
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required,base64"`
}

type MailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
