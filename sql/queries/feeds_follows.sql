-- name: CreateFeedFollow :one
WITH inserted AS (
  INSERT INTO feed_follows (user_id, feed_id)
  VALUES ($1, $2)
  RETURNING *
)
SELECT inserted.*, feeds.name as feed_name, users.name as user_name
  FROM inserted
  JOIN feeds ON inserted.feed_id = feeds.id
  JOIN users ON inserted.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name as feed_name
FROM feed_follows
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;
