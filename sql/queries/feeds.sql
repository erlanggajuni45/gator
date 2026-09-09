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
