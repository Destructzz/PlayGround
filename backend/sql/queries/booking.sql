-- name: CreateBooking :one
WITH service_vars AS (
    SELECT price / NULLIF(duration, 0) AS price_per_minute
    FROM services
    WHERE id = sqlc.arg(service_id)
)
INSERT INTO bookings (
    user_id,
    zone_id,
    service_id,
    start_time,
    end_time,
    participants,
    total_price,
    status,
    created_at,
    updated_at
)
SELECT
    sqlc.arg(user_id),
    sqlc.arg(zone_id),
    sqlc.arg(service_id),
    sqlc.arg(start_time),
    sqlc.arg(end_time),
    sqlc.arg(participants)::int,
    FLOOR(sv.price_per_minute * (EXTRACT(EPOCH FROM (sqlc.arg(end_time)::timestamptz - sqlc.arg(start_time)::timestamptz)) / 60) * sqlc.arg(participants)::numeric),
    sqlc.arg(status),
    NOW(),
    NOW()
FROM service_vars sv
RETURNING
    id,
    user_id,
    zone_id,
    service_id,
    start_time,
    end_time,
    participants,
    total_price,
    status,
    created_at,
    updated_at;

-- name: ListBookings :many
SELECT id, user_id, zone_id, service_id, start_time, end_time, participants, total_price, status, created_at, updated_at
FROM bookings
ORDER BY id;

-- name: GetBookingByID :one
SELECT id, user_id, zone_id, service_id, start_time, end_time, participants, total_price, status, created_at, updated_at
FROM bookings
WHERE id = $1;

-- name: PatchBooking :one
WITH current_booking AS (
    SELECT *
    FROM bookings
    WHERE id = sqlc.arg(id)
),
new_values AS (
    SELECT
        COALESCE(sqlc.narg(user_id), cb.user_id) AS user_id,
        COALESCE(sqlc.narg(zone_id), cb.zone_id) AS zone_id,
        COALESCE(sqlc.narg(service_id), cb.service_id) AS service_id,
        COALESCE(sqlc.narg(start_time), cb.start_time) AS start_time,
        COALESCE(sqlc.narg(end_time), cb.end_time) AS end_time,
        COALESCE(sqlc.narg(participants), cb.participants) AS participants,
        COALESCE(sqlc.narg(status), cb.status) AS status
    FROM current_booking cb
),
service_vars AS (
    SELECT
        s.id,
        s.price / NULLIF(s.duration, 0) AS price_per_minute
    FROM services s
    JOIN new_values nv ON nv.service_id = s.id
)
UPDATE bookings b
SET
    user_id = nv.user_id,
    zone_id = nv.zone_id,
    service_id = nv.service_id,
    start_time = nv.start_time,
    end_time = nv.end_time,
    participants = nv.participants,
    total_price = FLOOR(
        sv.price_per_minute *
        (EXTRACT(EPOCH FROM (nv.end_time - nv.start_time)) / 60) *
        nv.participants
    ),
    status = nv.status,
    updated_at = NOW()
FROM new_values nv
JOIN service_vars sv ON sv.id = nv.service_id
WHERE b.id = sqlc.arg(id)
RETURNING
    b.id,
    b.user_id,
    b.zone_id,
    b.service_id,
    b.start_time,
    b.end_time,
    b.participants,
    b.total_price,
    b.status,
    b.created_at,
    b.updated_at;

-- name: DeleteBooking :one
DELETE FROM bookings
WHERE id = $1
RETURNING id;
