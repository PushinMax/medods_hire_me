package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS users(
    		userID UUID primary key,
    		email text
		)`,
	); err != nil {
		return nil, err
	}

	if _, err := db.Exec(
		`CREATE TABLE IF not EXISTS tokens(
    		jti UUID primary key,
    		userID UUID not null,
    		ip text not null,
    		token_hash TEXT unique,
    		expires_at TIMESTAMP not null,

    		FOREIGN KEY (userID) REFERENCES users(userID)
        		ON DELETE CASCADE
		)`,
	); err != nil {
		return nil, err
	}

	return db, nil
}