-- name: ListPings :many
SELECT id, content, created_at
FROM pings
ORDER BY created_at DESC;

-- name: CreatePing :one
INSERT INTO pings (content)
VALUES ($1)
RETURNING id, content, created_at;
