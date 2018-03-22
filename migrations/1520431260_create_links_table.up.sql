CREATE TABLE links (
	hash varchar(64) NOT NULL PRIMARY KEY,
	url text NOT NULL,
	status int NOT NULL DEFAULT 1, -- 0 - Deleted; 1 - Active
	created_at timestamptz NOT NULL
);