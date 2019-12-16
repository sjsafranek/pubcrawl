

DROP TABLE IF EXISTS venues CASCADE;
CREATE TABLE venues (
	id			VARCHAR NOT NULL,
	crawl_id	VARCHAR,
	visited		BOOLEAN DEFAULT false,
	data		JSONB,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (crawl_id) REFERENCES crawls(id) ON DELETE CASCADE,
	PRIMARY KEY (id, crawl_id)
);
CREATE INDEX ON venues (crawl_id);

-- @trigger crawls_update
DROP TRIGGER IF EXISTS venues_update ON venues;
CREATE TRIGGER venues_update
    BEFORE UPDATE ON venues
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();
