
DROP VIEW IF EXISTS venues_view CASCADE;
CREATE OR REPLACE VIEW venues_view AS (

    -- WITH next AS (
    --     SELECT
    --         votes_view.crawl_id,
    --         votes_view.venue_id,
    --         votes_view.votes
    --     FROM votes_view
    --     WHERE votes_view.visited = false
    --     ORDER BY votes_view.votes DESC
    --     LIMIT 1
    -- )

    SELECT
        crawl_id,
        json_build_object(
            -- 'next', NULL,

            'next', (
                SELECT
                    -- votes_view.crawl_id,
                    json_agg(votes_view.venue_id)
                    -- votes_view.votes
                FROM votes_view
                WHERE votes_view.visited = false
                GROUP BY votes_view.crawl_id, votes_view.venue_id, votes_view.votes
                ORDER BY votes_view.votes DESC
                LIMIT 1
            ),

            'visited', (
                SELECT
                    json_agg(visited.data)
                FROM venues AS visited
                WHERE visited.visited = true
                AND  visited.crawl_id = venues.crawl_id
                GROUP BY visited.crawl_id
            ),
            'unvisited', (
                SELECT
                    json_agg(unvisited.data)
                FROM venues AS unvisited
                WHERE unvisited.visited = false
                AND  unvisited.crawl_id = venues.crawl_id
                GROUP BY unvisited.crawl_id
            )
        ) AS venues_json
    FROM venues
    GROUP BY venues.crawl_id
);
