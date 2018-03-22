CREATE TABLE users_x_links (
	iduser int NOT NULL REFERENCES users ON DELETE CASCADE,
	hashlink varchar(64) NOT NULL REFERENCES links ON DELETE CASCADE,
	description text NULL DEFAULT '',
	status int NOT NULL DEFAULT 1, -- 0 - Deleted; 1 - Active
	visibility int NOT NULL DEFAULT 1 -- 0 - Private; 1 - Public; 2 - Friends only
);

CREATE INDEX users_x_links_iduser_hashlink_idx ON users_x_links (iduser, hashlink);