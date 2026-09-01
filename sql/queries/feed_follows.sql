-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    inserted_feed_follow.id,
    inserted_feed_follow.created_at,
    inserted_feed_follow.updated_at,
    inserted_feed_follow.feed_id,
    inserted_feed_follow.user_id,
    feeds.name AS feed_name,
    users.name AS user_name
from inserted_feed_follow
INNER JOIN feeds
ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users
ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsPerUser :many
SELECT
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE users.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 and feed_id = $2;