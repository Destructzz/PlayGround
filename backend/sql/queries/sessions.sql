-- name: CreateSession :one
INSERT INTO sessions (user_id)
VALUES ($1)
RETURNING id, user_id, created_at, expires_at;

-- name: GetSessionByID :one
SELECT s.id, s.user_id, s.created_at, s.expires_at,
       u.id AS "user_id", u.google_id, u.full_name, u.email, u.avatar_url, u.phone, u.role
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.id = $1 AND s.expires_at > NOW();

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: GetUserByID :one
SELECT id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at
FROM users
WHERE id = $1;
