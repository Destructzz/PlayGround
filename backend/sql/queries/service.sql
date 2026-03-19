-- name: CreateService :one
INSERT INTO services (name, zone_id, duration, price, currency, description, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING id, name, zone_id, duration, price, currency, description, is_active, created_at, updated_at;

-- name: ListServices :many
SELECT id, name, zone_id, duration, price, currency, description, is_active, created_at, updated_at FROM services;

-- name: GetServiceByID :one
SELECT id, name, zone_id, duration, price, currency, description, is_active, created_at, updated_at FROM services WHERE id = $1;

-- name: PatchService :one
UPDATE services
SET name = COALESCE(sqlc.narg(name), name),
    zone_id = COALESCE(sqlc.narg(zone_id), zone_id),
    duration = COALESCE(sqlc.narg(duration), duration),
    price = COALESCE(sqlc.narg(price), price),
    currency = COALESCE(sqlc.narg(currency), currency),
    description = COALESCE(sqlc.narg(description), description),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING id, name, zone_id, duration, price, currency, description, is_active, created_at, updated_at;

-- name: DeleteService :one
DELETE FROM services WHERE id = $1 RETURNING id;
