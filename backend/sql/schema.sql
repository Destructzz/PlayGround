CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE role AS ENUM ('admin', 'seller', 'client');
CREATE TYPE zone_type AS ENUM ('game', 'event', 'lounge', 'sys');
CREATE TYPE booking_status AS ENUM ('created', 'confirmed', 'canceled', 'completed');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'online');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded');
CREATE TYPE action_type AS ENUM ('create', 'update', 'delete', 'login', 'logout');
CREATE TYPE entity_type AS ENUM ('booking', 'service', 'zone', 'payment', 'user', 'staff', 'other');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR,
    full_name TEXT NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    avatar_url VARCHAR,
    phone VARCHAR,
    role role NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '7 days'
);

CREATE TABLE IF NOT EXISTS zone_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zones (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    zone_type zone_type NOT NULL,
    zone_tag_id int REFERENCES zone_tags(id) ON DELETE CASCADE,
    capacity INTEGER NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    zone_id BIGINT NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    duration INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    currency VARCHAR NOT NULL DEFAULT 'RUB',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS computer_configurations(
    id BIGSERIAL PRIMARY KEY,
    zone_tags_id BIGINT NOT NULL REFERENCES zone_tags(id) ON DELETE CASCADE,
    specs_json JSONB NOT NULL DEFAULT '[]'::jsonb
);

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

CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    zone_id BIGINT NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    place_id BIGINT REFERENCES zone_places(id) ON DELETE SET NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    participants INTEGER NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL,
    status booking_status NOT NULL DEFAULT 'created',
    contact_name VARCHAR NOT NULL DEFAULT '',
    contact_email VARCHAR NOT NULL DEFAULT '',
    contact_phone VARCHAR NOT NULL DEFAULT '',
    details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS shifts (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    zone_tag_id BIGINT REFERENCES zone_tags(id) ON DELETE SET NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT shifts_valid_time_range CHECK (end_time > start_time)
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    amount NUMERIC(10, 2) NOT NULL,
    currency VARCHAR NOT NULL DEFAULT 'RUB',
    payment_method payment_method NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    receipt_number VARCHAR,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action action_type NOT NULL,
    entity entity_type NOT NULL,
    entity_id BIGINT,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR
);

CREATE TABLE IF NOT EXISTS pings (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS site_settings(
    id INT PRIMARY KEY DEFAULT 1 CHECK (id = 1), -- Singleton table (только 1 строка всегда)
    settings_json JSONB NOT NULL DEFAULT '[]'::jsonb, -- ID фич или их массив для отображения
    gallery_items_json JSONB NOT NULL DEFAULT '[]'::jsonb, -- То, что каруселится внизу
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
