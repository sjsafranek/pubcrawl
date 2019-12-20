
DROP VIEW IF EXISTS crawls_view CASCADE;
CREATE OR REPLACE VIEW crawls_view AS (
    SELECT
        *,
        json_build_object(
            'id', id,
            'name', name,
            'owner', owner,
            'max_votes_per_user', max_votes_per_user,
            'is_deleted', is_deleted,
            'created_at', to_char(created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
            'updated_at', to_char(updated_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"')
        ) AS crawls_json
    FROM crawls
	LEFT JOIN venues_view
		ON crawls.id = venues_view.crawl_id
    WHERE
        is_deleted = false
);
