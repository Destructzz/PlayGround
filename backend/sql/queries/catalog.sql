-- name: ListZoneTags :many
SELECT id, name, created_at, updated_at
FROM zone_tags
ORDER BY id;

-- name: ListUserZoneTags :many
SELECT id, name, created_at, updated_at
FROM zone_tags
ORDER BY id
LIMIT 3;

-- name: GetZoneTagByID :one
SELECT id, name, created_at, updated_at
FROM zone_tags
WHERE id = $1;

-- name: CreateZoneTag :one
INSERT INTO zone_tags (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: PatchZoneTag :one
UPDATE zone_tags
SET
    name = COALESCE(sqlc.narg(name), name),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, created_at, updated_at;

-- name: DeleteZoneTag :execrows
DELETE FROM zone_tags
WHERE id = $1;

-- name: GetSiteSettings :one
SELECT id, settings_json, gallery_items_json, updated_at
FROM site_settings
WHERE id = 1;

-- name: UpsertSiteSettings :one
INSERT INTO site_settings (id, settings_json, gallery_items_json, updated_at)
VALUES (1, $1, $2, NOW())
ON CONFLICT (id) DO UPDATE
SET settings_json = EXCLUDED.settings_json,
    gallery_items_json = EXCLUDED.gallery_items_json,
    updated_at = NOW()
RETURNING id, settings_json, gallery_items_json, updated_at;

-- name: ListComputerConfigurations :many
SELECT id, zone_tags_id, specs_json
FROM computer_configurations
ORDER BY id;

-- name: ListComputerConfigurationsByZoneTagIDs :many
SELECT id, zone_tags_id, specs_json
FROM computer_configurations
WHERE zone_tags_id = ANY($1::bigint[])
ORDER BY zone_tags_id, id;

-- name: GetComputerConfigurationByID :one
SELECT id, zone_tags_id, specs_json
FROM computer_configurations
WHERE id = $1;

-- name: CreateComputerConfiguration :one
INSERT INTO computer_configurations (zone_tags_id, specs_json)
VALUES ($1, $2)
RETURNING id, zone_tags_id, specs_json;

-- name: PatchComputerConfiguration :one
UPDATE computer_configurations
SET
    zone_tags_id = COALESCE(sqlc.narg(zone_tags_id), zone_tags_id),
    specs_json = COALESCE(sqlc.narg(specs_json), specs_json)
WHERE id = $1
RETURNING id, zone_tags_id, specs_json;

-- name: DeleteComputerConfiguration :execrows
DELETE FROM computer_configurations
WHERE id = $1;
