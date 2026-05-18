-- name: UpsertUser :one
INSERT INTO users (google_id, full_name, email, avatar_url, role)
VALUES ($1, $2, $3, $4, 'client')
ON CONFLICT (email) DO UPDATE
SET google_id = EXCLUDED.google_id,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = NOW()
RETURNING id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at;

-- name: PatchUser :one
UPDATE users
SET full_name = COALESCE(sqlc.narg('full_name'), full_name),
    phone = COALESCE(sqlc.narg('phone'), phone),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at;

-- name: ListUsers :many
SELECT id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY full_name;