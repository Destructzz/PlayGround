-- name: CreatePayment :one
INSERT INTO payments (booking_id, amount, currency, payment_method, status, receipt_number, paid_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, booking_id, amount, currency, payment_method, status, receipt_number, paid_at, created_at, updated_at;

-- name: ListPayments :many
SELECT id, booking_id, amount, currency, payment_method, status, receipt_number, paid_at, created_at, updated_at
FROM payments
ORDER BY id;

-- name: GetPaymentByID :one
SELECT id, booking_id, amount, currency, payment_method, status, receipt_number, paid_at, created_at, updated_at
FROM payments
WHERE id = $1;

-- name: PatchPayment :one
UPDATE payments
SET booking_id = COALESCE(sqlc.narg(booking_id), booking_id),
    amount = COALESCE(sqlc.narg(amount), amount),
    currency = COALESCE(sqlc.narg(currency), currency),
    payment_method = COALESCE(sqlc.narg(payment_method), payment_method),
    status = COALESCE(sqlc.narg(status), status),
    receipt_number = COALESCE(sqlc.narg(receipt_number), receipt_number),
    paid_at = COALESCE(sqlc.narg(paid_at), paid_at),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING id, booking_id, amount, currency, payment_method, status, receipt_number, paid_at, created_at, updated_at;

-- name: DeletePayment :one
DELETE FROM payments
WHERE id = $1
RETURNING id;
