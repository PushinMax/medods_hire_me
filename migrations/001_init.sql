CREATE TABLE IF NOT EXISTS users(
    userID UUID primary key,
    email text
);

CREATE TABLE IF not EXISTS tokens(
    jti UUID primary key,
    userID UUID not null,
    ip text not null,
    token_hash TEXT unique,
    expires_at TIMESTAMP not null,

    FOREIGN KEY (userID) REFERENCES users(userID)
        ON DELETE CASCADE
);

