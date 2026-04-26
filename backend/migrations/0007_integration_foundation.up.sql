-- Migration 0007: Integration foundation
-- Adds session support, metadata JSONB fields, zone_places, booking enhancements

-- 1. Sessions table for cookie-based browser auth
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '7 days'
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

-- 2. Add details_json metadata to zones and services
ALTER TABLE zones ADD COLUMN IF NOT EXISTS details_json JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE services ADD COLUMN IF NOT EXISTS details_json JSONB NOT NULL DEFAULT '{}'::jsonb;

-- 3. Add zone_places table for gaming seats/places
CREATE TABLE IF NOT EXISTS zone_places (
    id BIGSERIAL PRIMARY KEY,
    zone_id BIGINT NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    label VARCHAR NOT NULL,
    configuration_id BIGINT REFERENCES computer_configurations(id) ON DELETE SET NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_zone_places_zone_id ON zone_places(zone_id);

-- 4. Extend bookings with place_id, contact snapshot, and details
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS place_id BIGINT REFERENCES zone_places(id) ON DELETE SET NULL;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS contact_name VARCHAR NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS contact_email VARCHAR NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS contact_phone VARCHAR NOT NULL DEFAULT '';
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS details_json JSONB NOT NULL DEFAULT '{}'::jsonb;

-- 5. Index for overlap detection on bookings
CREATE INDEX IF NOT EXISTS idx_bookings_overlap ON bookings(zone_id, service_id, start_time, end_time)
    WHERE status NOT IN ('canceled');
CREATE INDEX IF NOT EXISTS idx_bookings_place_overlap ON bookings(place_id, start_time, end_time)
    WHERE place_id IS NOT NULL AND status NOT IN ('canceled');
