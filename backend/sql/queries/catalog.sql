-- name: ListZoneTags :many
SELECT id, name, created_at, updated_at
FROM zone_tags
ORDER BY id;

-- name: GetZoneTagByID :one
SELECT id, name, created_at, updated_at
FROM zone_tags
WHERE id = $1;

-- name: CreateZoneTag :one
INSERT INTO zone_tags (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

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

-- name: GetComputerConfigurationByID :one
SELECT id, zone_tags_id, specs_json
FROM computer_configurations
WHERE id = $1;

-- name: CreateComputerConfiguration :one
INSERT INTO computer_configurations (zone_tags_id, specs_json)
VALUES ($1, $2)
RETURNING id, zone_tags_id, specs_json;
