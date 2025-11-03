-- +goose Up
-- Schema inicial: assets, work_orders, maintenance_plans

CREATE TABLE IF NOT EXISTS assets (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL,
    location       TEXT,
    criticality    TEXT CHECK (criticality IN ('A','B','C')) DEFAULT 'B',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_orders (
    id                BIGSERIAL PRIMARY KEY,
    asset_id          BIGINT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    type              TEXT NOT NULL CHECK (type IN ('corrective','preventive','condition','improvement')) DEFAULT 'corrective',
    status            TEXT NOT NULL CHECK (status IN ('open','in_progress','done','canceled')) DEFAULT 'open',
    title             TEXT NOT NULL,
    description       TEXT,
    breakdown_at      TIMESTAMPTZ,
    closed_at         TIMESTAMPTZ,
    downtime_minutes  BIGINT,
    cause             TEXT,
    solution          TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS maintenance_plans (
    id              BIGSERIAL PRIMARY KEY,
    asset_id        BIGINT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    rule_type       TEXT NOT NULL CHECK (rule_type IN ('time','meter','condition')),
    frequency_days  BIGINT,
    meter_target    BIGINT,
    last_execution  TIMESTAMPTZ,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_time_requires_freq CHECK (
        (rule_type <> 'time') OR (frequency_days IS NOT NULL)
    )
);

-- +goose Down
DROP TABLE IF EXISTS maintenance_plans;
DROP TABLE IF EXISTS work_orders;
DROP TABLE IF EXISTS assets;
