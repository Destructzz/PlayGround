-- name: CreateShift :one
INSERT INTO shifts (user_id, zone_tag_id, start_time, end_time, note)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, zone_tag_id, start_time, end_time, note, created_at, updated_at, deleted_at;

-- name: ListShifts :many
SELECT sqlc.embed(shifts), sqlc.embed(users)
FROM shifts
JOIN users ON users.id = shifts.user_id
WHERE shifts.deleted_at IS NULL
ORDER BY shifts.start_time DESC, shifts.id DESC;

-- name: GetShiftByID :one
SELECT sqlc.embed(shifts), sqlc.embed(users)
FROM shifts
JOIN users ON users.id = shifts.user_id
WHERE shifts.id = $1 AND shifts.deleted_at IS NULL;

-- name: PatchShift :one
UPDATE shifts
SET zone_tag_id = COALESCE(sqlc.narg(zone_tag_id), zone_tag_id),
    start_time = COALESCE(sqlc.narg(start_time), start_time),
    end_time = COALESCE(sqlc.narg(end_time), end_time),
    note = COALESCE(sqlc.narg(note), note),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id, user_id, zone_tag_id, start_time, end_time, note, created_at, updated_at, deleted_at;

-- name: DeleteShift :one
UPDATE shifts
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;

-- name: ListShiftsByZoneTagID :many
SELECT sqlc.embed(shifts), sqlc.embed(users)
FROM shifts
JOIN users ON users.id = shifts.user_id
WHERE shifts.zone_tag_id = $1
  AND shifts.deleted_at IS NULL
ORDER BY shifts.start_time ASC;

-- name: HasOverlappingShift :one
-- Returns true if any existing shift for the same zone_tag_id overlaps with the given window.
-- exclude_id is used on PATCH to exclude the shift being updated (pass 0 to skip).
SELECT EXISTS (
  SELECT 1
  FROM shifts
  WHERE zone_tag_id = sqlc.arg(zone_tag_id)
    AND deleted_at IS NULL
    AND id <> sqlc.arg(exclude_id)
    AND start_time < sqlc.arg(end_time)
    AND end_time   > sqlc.arg(start_time)
) AS overlaps;
