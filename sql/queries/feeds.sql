-- name: CreateFeed :one
INSERT INTO feeds (
    name,
    url,
    user_id
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetFeeds :many
SELECT a.name, a.url, b.name AS user FROM feeds a
JOIN users b ON a.user_id = b.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds SET last_fetched_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;
