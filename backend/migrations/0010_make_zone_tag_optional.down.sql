BEGIN;

-- We can't easily revert if there are NULLs, so we assign a default tag if reverting.
-- Assuming 'default' tag exists from migration 0005.
UPDATE zones SET zone_tag_id = (SELECT id FROM zone_tags WHERE name = 'default' LIMIT 1) WHERE zone_tag_id IS NULL;

ALTER TABLE zones ALTER COLUMN zone_tag_id SET NOT NULL;

COMMIT;
