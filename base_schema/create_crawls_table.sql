
DROP TABLE IF EXISTS crawls CASCADE;
CREATE TABLE crawls (
	id					VARCHAR PRIMARY KEY NOT NULL DEFAULT md5(random()::text||now()::text)::uuid,
    name                VARCHAR,
	owner 				VARCHAR,
	max_votes_per_user	INTEGER DEFAULT 1,
    is_deleted          BOOLEAN DEFAULT false,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (owner) REFERENCES users(username) ON DELETE CASCADE
);

-- @trigger crawls_update
DROP TRIGGER IF EXISTS crawls_update ON crawls;
CREATE TRIGGER crawls_update
    BEFORE UPDATE ON crawls
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();
