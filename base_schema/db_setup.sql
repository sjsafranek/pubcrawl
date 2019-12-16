/*======================================================================*/
--  db_setup.sql
--   -- :mode=pl-sql:tabSize=3:indentSize=3:
--  Mon Aug 17 14:44:44 PST 2015 @144 /Internet Time/
--  Purpose:
--  NOTE: must be connected as 'postgres' user or a superuser to start.
/*======================================================================*/

\set ON_ERROR_STOP on
set client_min_messages to 'warning';


\i create_extensions.sql
\i create_general_functions.sql
\i create_config_table.sql
\i create_users_table.sql
\i create_social_accounts_table.sql
\i create_crawls_table.sql
\i create_venues_table.sql









DROP TABLE IF EXISTS up_votes CASCADE;
CREATE TABLE up_votes (
	owner		VARCHAR,
	crawl_id 	VARCHAR,
	venue_id	VARCHAR,
	FOREIGN KEY (owner) REFERENCES users(username) ON DELETE CASCADE,
	FOREIGN KEY (crawl_id) REFERENCES crawls(id) ON DELETE CASCADE,
	FOREIGN KEY (venue_id, crawl_id) REFERENCES venues(id, crawl_id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS down_votes CASCADE;
CREATE TABLE down_votes (
	owner 		VARCHAR,
	crawl_id 	VARCHAR,
	venue_id	VARCHAR,
	FOREIGN KEY (owner) REFERENCES users(username) ON DELETE CASCADE,
	FOREIGN KEY (crawl_id) REFERENCES crawls(id) ON DELETE CASCADE,
	FOREIGN KEY (venue_id, crawl_id) REFERENCES venues(id, crawl_id) ON DELETE CASCADE
);

-- DROP TABLE IF EXISTS user_voting_rules;
-- CREATE TABLE user_voting_rules (
-- 	owner 				VARCHAR,
-- 	crawl_id 			VARCHAR,
-- 	up_vote_weight		INTEGER DEFAULT 1,
-- 	down_vote_weight	INTEGER DEFAULT 1,
-- 	FOREIGN KEY (owner) REFERENCES users(username) ON DELETE CASCADE,
-- 	FOREIGN KEY (crawl_id) REFERENCES crawls(id) ON DELETE CASCADE
-- );

DROP VIEW IF EXISTS votes_view CASCADE;
CREATE OR REPLACE VIEW votes_view AS (
	WITH upv AS (
		SELECT crawl_id, venue_id, COUNT(*) AS votes FROM up_votes GROUP BY crawl_id, venue_id
	),
		dnv AS (
		SELECT crawl_id, venue_id, COUNT(*) AS votes FROM down_votes GROUP BY crawl_id, venue_id
	)
	SELECT
		venues.visited AS visited,
		venues.crawl_id AS crawl_id,
		venues.id AS venue_id,
		upv.votes - dnv.votes AS votes
	FROM venues
	LEFT JOIN upv
		ON upv.crawl_id = venues.crawl_id
		AND upv.venue_id = venues.id
	LEFT JOIN dnv
		ON dnv.crawl_id = venues.crawl_id
		AND dnv.venue_id = venues.id
);








\i create_users_view.sql
\i create_venues_view.sql
\i create_crawls_view.sql
