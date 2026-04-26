-- Rollback migration 0007
DROP INDEX IF EXISTS idx_bookings_place_overlap;
DROP INDEX IF EXISTS idx_bookings_overlap;

ALTER TABLE bookings DROP COLUMN IF EXISTS details_json;
ALTER TABLE bookings DROP COLUMN IF EXISTS contact_phone;
ALTER TABLE bookings DROP COLUMN IF EXISTS contact_email;
ALTER TABLE bookings DROP COLUMN IF EXISTS contact_name;
ALTER TABLE bookings DROP COLUMN IF EXISTS place_id;

DROP TABLE IF EXISTS zone_places;

ALTER TABLE services DROP COLUMN IF EXISTS details_json;
ALTER TABLE zones DROP COLUMN IF EXISTS details_json;

DROP TABLE IF EXISTS sessions;
