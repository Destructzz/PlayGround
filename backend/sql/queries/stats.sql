-- name: GetSummaryStats :one
SELECT
    (SELECT COUNT(*) FROM bookings) AS total_bookings,
    (SELECT COALESCE(SUM(amount), 0)::NUMERIC FROM payments WHERE status = 'paid') AS total_revenue,
    (SELECT COUNT(*) FROM users WHERE is_active = TRUE) AS active_users,
    (SELECT COUNT(*) FROM bookings WHERE status = 'created') AS pending_bookings;

-- name: GetRevenueLast30Days :many
SELECT
    DATE_TRUNC('day', paid_at)::DATE AS date,
    SUM(amount)::NUMERIC AS revenue
FROM payments
WHERE status = 'paid' AND paid_at > NOW() - INTERVAL '30 days'
GROUP BY 1
ORDER BY 1;

-- name: GetBookingsLast30Days :many
SELECT
    DATE_TRUNC('day', created_at)::DATE AS date,
    COUNT(*) AS count
FROM bookings
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY 1
ORDER BY 1;

-- name: GetBookingsByZoneType :many
SELECT
    z.zone_type,
    COUNT(*) AS count
FROM bookings b
JOIN zones z ON b.zone_id = z.id
GROUP BY z.zone_type;
