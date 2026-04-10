BEGIN;

-- 1. Remove constraints and column from zones
ALTER TABLE zones DROP CONSTRAINT IF EXISTS fk_zone_tags;
ALTER TABLE zones DROP COLUMN IF EXISTS zone_tag_id;

-- 2. Drop table zone_tags
DROP TABLE IF EXISTS zone_tags;

-- 3. Restore 'vip' to enum
ALTER TYPE zone_type RENAME TO zone_type_old;
CREATE TYPE zone_type AS ENUM ('game', 'event', 'vip', 'lounge', 'sys');
ALTER TABLE zones ALTER COLUMN zone_type TYPE zone_type USING zone_type::text::zone_type;
DROP TYPE zone_type_old;

COMMIT;
