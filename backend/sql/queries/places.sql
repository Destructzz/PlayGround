-- name: ListZonePlaces :many
SELECT id, zone_id, label, configuration_id, sort_order, is_active, created_at, updated_at
FROM zone_places
WHERE zone_id = $1
ORDER BY sort_order, id;

-- name: GetZonePlaceByID :one
SELECT id, zone_id, label, configuration_id, sort_order, is_active, created_at, updated_at
FROM zone_places
WHERE id = $1;

-- name: CreateZonePlace :one
INSERT INTO zone_places (zone_id, label, configuration_id, sort_order, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, zone_id, label, configuration_id, sort_order, is_active, created_at, updated_at;

-- name: DeleteZonePlace :execrows
DELETE FROM zone_places
WHERE id = $1;

-- name: UpdateZonePlace :one
UPDATE zone_places
SET
  label = COALESCE(sqlc.narg(label), label),
  configuration_id = COALESCE(sqlc.narg(configuration_id), configuration_id),
  sort_order = COALESCE(sqlc.narg(sort_order), sort_order),
  is_active = COALESCE(sqlc.narg(is_active), is_active),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING id, zone_id, label, configuration_id, sort_order, is_active, created_at, updated_at;

-- name: GetAllZonePlaces :many
SELECT id, zone_id, label, configuration_id, sort_order, is_active, created_at, updated_at
FROM zone_places
ORDER BY sort_order, id;


-- name: ListActiveZonePlaces :many
SELECT zp.id, zp.zone_id, zp.label, zp.configuration_id, zp.sort_order, zp.is_active, zp.created_at, zp.updated_at
FROM zone_places zp
WHERE zp.zone_id = $1 AND zp.is_active = TRUE
ORDER BY zp.sort_order, zp.id;

-- name: CheckPlaceBookingOverlap :one
SELECT COUNT(*) AS conflict_count
FROM bookings
WHERE place_id = $1
  AND status NOT IN ('canceled')
  AND start_time < sqlc.arg(end_time)
  AND end_time > sqlc.arg(start_time);

-- name: CheckZoneBookingOverlap :one
SELECT COUNT(*) AS conflict_count
FROM bookings
WHERE zone_id = $1
  AND service_id = $2
  AND status NOT IN ('canceled')
  AND start_time < sqlc.arg(end_time)
  AND end_time > sqlc.arg(start_time);

-- name: CountActiveBookingsForZone :one
SELECT COALESCE(SUM(participants), 0)::int AS booked_count
FROM bookings
WHERE zone_id = $1
  AND status NOT IN ('canceled')
  AND start_time < sqlc.arg(end_time)
  AND end_time > sqlc.arg(start_time);

-- name: ListZonePlaceBookingsForDate :many
SELECT id, zone_id, service_id, place_id, start_time, end_time, status
FROM bookings
WHERE zone_id = sqlc.arg(zone_id)
  AND place_id IS NOT NULL
  AND status NOT IN ('canceled')
  AND start_time < sqlc.arg(date_end)
  AND end_time > sqlc.arg(date_start)
ORDER BY start_time, id;

-- name: ListLoungeBookingsForDate :many
SELECT id, zone_id, start_time, end_time, participants, status
FROM bookings
WHERE zone_id = sqlc.arg(zone_id)
  AND status NOT IN ('canceled')
  AND start_time < sqlc.arg(date_end)
  AND end_time > sqlc.arg(date_start)
ORDER BY start_time, id;
