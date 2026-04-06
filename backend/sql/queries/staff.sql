-- name: CreateStaff :one
INSERT INTO staff (user_id, position, hire_date, phone, email, is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, position, hire_date, phone, email, is_active, created_at, updated_at, deleted_at;

-- name: ListStaff :many
SELECT id, user_id, position, hire_date, phone, email, is_active, created_at, updated_at, deleted_at
FROM staff
WHERE deleted_at IS NULL
ORDER BY id;

-- name: GetStaffByID :one
SELECT id, user_id, position, hire_date, phone, email, is_active, created_at, updated_at, deleted_at
FROM staff
WHERE id = $1 AND deleted_at IS NULL;

-- name: PatchStaff :one
UPDATE staff
SET position = COALESCE(sqlc.narg(position), position),
    hire_date = COALESCE(sqlc.narg(hire_date), hire_date),
    phone = COALESCE(sqlc.narg(phone), phone),
    email = COALESCE(sqlc.narg(email), email),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id, user_id, position, hire_date, phone, email, is_active, created_at, updated_at, deleted_at;

-- name: DeleteStaff :one
UPDATE staff
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;
