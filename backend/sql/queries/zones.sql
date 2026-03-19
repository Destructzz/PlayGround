-- name: CreateZone :one
INSERT INTO zones (name, zone_type, capacity, description, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, zone_type, capacity, description, is_active, created_at, updated_at;

-- name: ListZones :many
SELECT id, name, zone_type, capacity, description, is_active, created_at, updated_at
FROM zones
ORDER BY id;

-- name: GetZoneByID :one
SELECT id, name, zone_type, capacity, description, is_active, created_at, updated_at
FROM zones
WHERE id = $1;

-- name: PatchZone :one
UPDATE zones
SET name = COALESCE(sqlc.narg(name), name),
    zone_type = COALESCE(sqlc.narg(zone_type), zone_type),
    capacity = COALESCE(sqlc.narg(capacity), capacity),
    description = COALESCE(sqlc.narg(description), description),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING id, name, zone_type, capacity, description, is_active, created_at, updated_at;

-- name: DeleteZone :one
DELETE FROM zones
WHERE id = $1
RETURNING id;
