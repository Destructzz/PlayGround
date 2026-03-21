-- name: UpsertUser :one
INSERT INTO users (google_id, full_name, email, avatar_url, role)
VALUES ($1, $2, $3, $4, 'client')
ON CONFLICT (email) DO UPDATE
SET google_id = EXCLUDED.google_id,
    full_name = EXCLUDED.full_name,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = NOW()
RETURNING id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at;