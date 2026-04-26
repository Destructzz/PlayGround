-- name: CreateZone :one
INSERT INTO zones (name, zone_type, zone_tag_id, capacity, description, is_active, details_json)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, zone_type, zone_tag_id, capacity, description, is_active, details_json, created_at, updated_at;

-- name: ListZones :many
SELECT id, name, zone_type, zone_tag_id, capacity, description, is_active, details_json, created_at, updated_at
FROM zones
ORDER BY id;

-- name: GetZoneByID :one
SELECT id, name, zone_type, zone_tag_id, capacity, description, is_active, details_json, created_at, updated_at
FROM zones
WHERE id = $1;

-- name: PatchZone :one
UPDATE zones
SET name = COALESCE(sqlc.narg(name), name),
    zone_type = COALESCE(sqlc.narg(zone_type), zone_type),
    zone_tag_id = COALESCE(sqlc.narg(zone_tag_id), zone_tag_id),
    capacity = COALESCE(sqlc.narg(capacity), capacity),
    description = COALESCE(sqlc.narg(description), description),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    details_json = COALESCE(sqlc.narg(details_json), details_json),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING id, name, zone_type, zone_tag_id, capacity, description, is_active, details_json, created_at, updated_at;

-- name: DeleteZone :one
DELETE FROM zones
WHERE id = $1
RETURNING id;

-- name: ListZonesByType :many
SELECT id, name, zone_type, zone_tag_id, capacity, description, is_active, details_json, created_at, updated_at
FROM zones
WHERE zone_type = $1 AND is_active = TRUE
ORDER BY id;

-- name: ListZonesWithServices :many
SELECT z.id, z.name, z.zone_type, z.zone_tag_id, z.capacity, z.description, z.is_active, z.details_json, z.created_at, z.updated_at
FROM zones z
WHERE z.is_active = TRUE
ORDER BY z.id;
