-- name: CreateZone :one
INSERT INTO zones (name, zone_type, capacity, description)
VALUES ($1, $2, $3, $4)
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
SET name = COALESCE($2, name),
    zone_type = COALESCE($3, zone_type),
    capacity = COALESCE($4, capacity),
    description = COALESCE($5, description),
    is_active = COALESCE($6, is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, zone_type, capacity, description, is_active, created_at, updated_at;
