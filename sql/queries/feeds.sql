-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT
    feeds.name,
    feeds.url,
    users.name as username
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;


-- name: GetFeedID :one
SELECT id from feeds where url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $3 
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST;