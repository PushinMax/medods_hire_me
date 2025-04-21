package repository

import (
	//"hash"

	"medods_hire_me/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepoApi interface {
	GetEmail(userID string) (string, error)
	SetEmail(userID, email string) error
	Exists(userID string) error
	SaveRefreshToken(userID, hashedRefresh, ip, jti string, t time.Time) error
	TokenVerification(data *TokenData) error
}

type Repository struct {
	RepoApi
}

func New(db *sqlx.DB) *Repository {
	return &Repository{RepoApi: newPostgresImpl(db)}
}

type TokenData struct {
	RefreshToken string
	AccessToken  *utils.CustomClaims
	IP           string
}

type Token struct {
	Jti       uuid.UUID `db:"jti"`        // primary key
	UserID    uuid.UUID `db:"userid"`     // foreign key to users table
	IP        string    `db:"ip"`         // IP address
	TokenHash string    `db:"token_hash"` // unique token hash
	ExpiresAt time.Time `db:"expires_at"` // expiration timestamp
}
