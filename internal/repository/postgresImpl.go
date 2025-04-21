package repository

import (
	"fmt"
	"time"

	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type postgresImpl struct {
	db *sqlx.DB
}

func newPostgresImpl(db *sqlx.DB) *postgresImpl {
	return &postgresImpl{db: db}
}

func (r *postgresImpl) GetEmail(userID string) (string, error) {
	var str string
	err := r.db.Get(&str, "select email from users where userID = $1", &userID)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (r *postgresImpl) SetEmail(userID, email string) error {
	_, err := r.db.Exec(`Update users SET email = $1 WHERE userID = $2`, &email, &userID)
	return err 
}

func (r *postgresImpl) Exists(userID string) error {
	var user string
	if err := r.db.Get(&user, "select userID from users where userID = $1", &userID); err != nil {
		if _, err = r.db.Query(`Insert into users values($1, $2)`, &userID, "writeme@gmail.com"); err != nil {
			return fmt.Errorf("Exists: %s", err.Error())
		}
	}
	return nil
}
func (r *postgresImpl) SaveRefreshToken(userID, hashedRefresh, ip, jti string, t time.Time) error {
	if _, err := r.db.Exec(
		"INSERT INTO tokens (jti, userID, ip, token_hash, expires_at) VALUES ($1, $2, $3, $4, $5)",
		&jti, &userID, &ip, &hashedRefresh, &t,
	); err != nil {
		return fmt.Errorf("bad save new token: %s", err.Error())
	}
	return nil
}

func (r *postgresImpl) findToken(token string) (*Token, error) {
	var hashed []string
	if err := r.db.Select(&hashed, "select token_hash from tokens"); err != nil {
		return nil, err
	}
	for _, h := range hashed {
		if err := bcrypt.CompareHashAndPassword([]byte(h), []byte(token)); err == nil {
			var t Token
			if err := r.db.Get(
				&t,
				"SELECT * FROM tokens WHERE token_hash = $1",
				&h); err != nil {
				return nil, err
			}
			return &t, nil
		}
	}
	return nil, errors.New("unknown refresh token")
}
func (r *postgresImpl) TokenVerification(data *TokenData) error {
	token, err := r.findToken(data.RefreshToken)
	if err != nil {
		return fmt.Errorf("findToken: %s", err.Error())
	}

	if data.AccessToken.Subject != token.UserID.String() {
		return errors.New("different owners of tokens")
	}
	if (data.IP != data.AccessToken.IP) || (data.IP != token.IP) || (data.AccessToken.IP != token.IP) {
		return errors.New("IP do not match")
	}
	if data.AccessToken.JTI != token.Jti.String() {
		return errors.New("SEND MESSAGE: was used different version of tokens")
	}
	if time.Now().Add(3 *time.Hour).After(token.ExpiresAt) {
		return errors.New("refresh token has expired")
	}


	if _, err := r.db.Query(
		`DELETE from tokens WHERE jti = $1`,
		data.AccessToken.JTI,
	); err != nil {
		return fmt.Errorf("revoked tokes: %w", err)
	}
	return nil
}
