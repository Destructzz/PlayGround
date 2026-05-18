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
    place_id,
    start_time,
    end_time,
    participants,
    total_price,
    status,
    contact_name,
    contact_email,
    contact_phone,
    details_json,
    created_at,
    updated_at
)
SELECT
    sqlc.arg(user_id),
    sqlc.arg(zone_id),
    sqlc.arg(service_id),
    sqlc.narg(place_id),
    sqlc.arg(start_time),
    sqlc.arg(end_time),
    sqlc.arg(participants)::int,
    FLOOR(sv.price_per_minute * (EXTRACT(EPOCH FROM (sqlc.arg(end_time)::timestamptz - sqlc.arg(start_time)::timestamptz)) / 60) * sqlc.arg(participants)::numeric),
    sqlc.arg(status),
    sqlc.arg(contact_name),
    sqlc.arg(contact_email),
    sqlc.arg(contact_phone),
    sqlc.arg(details_json),
    NOW(),
    NOW()
FROM service_vars sv
RETURNING
    id,
    user_id,
    zone_id,
    service_id,
    place_id,
    start_time,
    end_time,
    participants,
    total_price,
    status,
    contact_name,
    contact_email,
    contact_phone,
    details_json,
    created_at,
    updated_at;

-- name: ListBookings :many
SELECT id, user_id, zone_id, service_id, place_id, start_time, end_time, participants, total_price, status, contact_name, contact_email, contact_phone, details_json, created_at, updated_at
FROM bookings
ORDER BY id;

-- name: GetBookingByID :one
SELECT id, user_id, zone_id, service_id, place_id, start_time, end_time, participants, total_price, status, contact_name, contact_email, contact_phone, details_json, created_at, updated_at
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
        COALESCE(sqlc.narg(status), cb.status) AS status,
        COALESCE(sqlc.narg(contact_name), cb.contact_name) AS contact_name,
        COALESCE(sqlc.narg(contact_email), cb.contact_email) AS contact_email,
        COALESCE(sqlc.narg(contact_phone), cb.contact_phone) AS contact_phone
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
    contact_name = nv.contact_name,
    contact_email = nv.contact_email,
    contact_phone = nv.contact_phone,
    updated_at = NOW()
FROM new_values nv
JOIN service_vars sv ON sv.id = nv.service_id
WHERE b.id = sqlc.arg(id)
RETURNING
    b.id,
    b.user_id,
    b.zone_id,
    b.service_id,
    b.place_id,
    b.start_time,
    b.end_time,
    b.participants,
    b.total_price,
    b.status,
    b.contact_name,
    b.contact_email,
    b.contact_phone,
    b.details_json,
    b.created_at,
    b.updated_at;

-- name: DeleteBooking :one
DELETE FROM bookings
WHERE id = $1
RETURNING id;

-- name: ListBookingsByUser :many
SELECT id, user_id, zone_id, service_id, place_id, start_time, end_time, participants, total_price, status, contact_name, contact_email, contact_phone, details_json, created_at, updated_at
FROM bookings
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListServicesByZoneID :many
SELECT id, name, zone_id, duration, price, currency, description, is_active, details_json, created_at, updated_at
FROM services
WHERE zone_id = $1 AND is_active = TRUE
ORDER BY id;