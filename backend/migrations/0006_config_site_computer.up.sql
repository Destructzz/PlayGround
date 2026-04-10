BEGIN;

CREATE TABLE IF NOT EXISTS computer_configurations(
    id BIGSERIAL PRIMARY KEY,
    zone_tags_id BIGINT NOT NULL REFERENCES zone_tags(id) ON DELETE CASCADE,
    specs_json JSONB NOT NULL DEFAULT '[]'::jsonb
);

CREATE TABLE IF NOT EXISTS site_settings(
    id INT PRIMARY KEY DEFAULT 1 CHECK (id = 1), -- Singleton table (только 1 строка всегда)
    settings_json JSONB NOT NULL DEFAULT '[]'::jsonb, -- ID фич или их массив для отображения
    gallery_items_json JSONB NOT NULL DEFAULT '[]'::jsonb, -- То, что каруселится внизу
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;