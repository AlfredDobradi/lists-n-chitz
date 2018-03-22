CREATE TABLE user_tokens (
    iduser int NOT NULL REFERENCES users ON DELETE CASCADE,
    token varchar(64) NOT NULL UNIQUE,
    expires TIMESTAMPTZ NOT NULL,
    address varchar(45) NOT NULL
)