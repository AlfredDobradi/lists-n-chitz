CREATE TABLE users (
	id serial PRIMARY KEY,
	email varchar(255) NOT NULL UNIQUE,
	password varchar(64) NOT NULL,
	created_at timestamptz NOT NULL,
	status smallint NOT NULL DEFAULT 1 -- 0 - Deleted; 1 - Active
);