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

-- name: SetUserRole :one
UPDATE users
SET role = sqlc.arg(role)::role,
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at;

-- name: SearchUsersByEmail :many
SELECT id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at
FROM users
WHERE email ILIKE '%' || sqlc.arg(query) || '%' AND deleted_at IS NULL
ORDER BY full_name
LIMIT 20;

-- name: ListSellers :many
SELECT id, google_id, full_name, email, avatar_url, phone, role, is_active, created_at, updated_at, deleted_at
FROM users
WHERE role = 'seller' AND deleted_at IS NULL
ORDER BY full_name;

-- name: GetBookingByIDPublic :one
SELECT id, user_id, zone_id, service_id, place_id, start_time, end_time, participants, total_price, status, contact_name, contact_email, contact_phone, details_json, created_at, updated_at
FROM bookings
WHERE id = $1;