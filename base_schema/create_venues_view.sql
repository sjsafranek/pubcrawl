\
DROP VIEW IF EXISTS venues_view CASCADE;
CREATE OR REPLACE VIEW venues_view AS (
    SELECT
        crawls.id AS crawl_id,
        (
            SELECT
                json_agg(
                    json_build_object(
                        'venue_id', venues.id,
                        'visited', venues.visited,
                        'up_votes', COALESCE(
                            (
                                SELECT
                                    SUM(1)
                                FROM up_votes AS upv
                                    WHERE  upv.crawl_id = venues.crawl_id
                                    AND upv.venue_id = venues.id
                            ),
                            0
                        ),
                        'down_votes', COALESCE(
                            (
                                SELECT
                                    SUM(1)
                                FROM down_votes AS dnv
                                    WHERE  dnv.crawl_id = venues.crawl_id
                                    AND dnv.venue_id = venues.id
                            ),
                            0
                        ),
                        'created_at', to_char(venues.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
                        'updated_at', to_char(venues.updated_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
                        'venue', venues.data
                    )
                )
            FROM venues
            WHERE venues.crawl_id = crawls.id
        ) AS venues_json
    FROM crawls
);





-- DROP VIEW IF EXISTS venues_view CASCADE;
-- CREATE OR REPLACE VIEW venues_view AS (
--     SELECT
--         crawl_id,
--         json_build_object(
--             'next', (
--                 SELECT
--                     json_agg(votes_view.venue_id)
--                 FROM votes_view
--                 WHERE votes_view.visited = false
--                 GROUP BY votes_view.crawl_id, votes_view.venue_id, votes_view.votes
--                 ORDER BY votes_view.votes DESC
--                 LIMIT 1
--             ),
--             'visited', (
--                 SELECT
--                     json_agg(visited.data)
--                 FROM venues AS visited
--                 WHERE visited.visited = true
--                 AND  visited.crawl_id = venues.crawl_id
--                 GROUP BY visited.crawl_id
--             ),
--             'unvisited', (
--                 SELECT
--                     json_agg(unvisited.data)
--                 FROM venues AS unvisited
--                 WHERE unvisited.visited = false
--                 AND  unvisited.crawl_id = venues.crawl_id
--                 GROUP BY unvisited.crawl_id
--             ),
--             'votes', (
--                 SELECT
--                     json_agg(votes_view.jdata)
--                 FROM votes_view
--                 WHERE
--                     votes_view.crawl_id = venues.crawl_id
--             )
--         ) AS venues_json
--     FROM venues
--     GROUP BY venues.crawl_id
-- );
