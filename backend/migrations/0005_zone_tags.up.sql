BEGIN;

-- 1. Modify enum to remove 'vip'
DELETE FROM zones WHERE zone_type = 'vip';
ALTER TYPE zone_type RENAME TO zone_type_old;
CREATE TYPE zone_type AS ENUM ('game', 'event', 'lounge', 'sys');
ALTER TABLE zones ALTER COLUMN zone_type TYPE zone_type USING zone_type::text::zone_type;
DROP TYPE zone_type_old;

-- 2. Create zone_tags
CREATE TABLE IF NOT EXISTS zone_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Add column to zones
-- We add it as nullable first, to avoid errors if there are existing rows
ALTER TABLE zones ADD COLUMN zone_tag_id int;

-- Insert a default tag to satisfy NOT NULL constraint for existing rows
INSERT INTO zone_tags (name) VALUES ('default') ON CONFLICT DO NOTHING;

-- Assign existing zones to the default tag
UPDATE zones SET zone_tag_id = (SELECT id FROM zone_tags WHERE name = 'default' LIMIT 1) WHERE zone_tag_id IS NULL;

-- 4. Apply NOT NULL and Foreign Key constraints
ALTER TABLE zones ALTER COLUMN zone_tag_id SET NOT NULL;
ALTER TABLE zones ADD CONSTRAINT fk_zone_tags FOREIGN KEY (zone_tag_id) REFERENCES zone_tags(id) ON DELETE CASCADE;

COMMIT;
