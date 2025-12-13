CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS plans (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price_cents INTEGER NOT NULL,
    billing_period TEXT NOT NULL,
    max_seats INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL UNIQUE,
    plan_id TEXT NOT NULL REFERENCES plans(id),
    seats INTEGER NOT NULL,
    status TEXT NOT NULL,
    activated_at TIMESTAMPTZ NOT NULL,
    current_period_end TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO plans (id, name, description, price_cents, billing_period, max_seats)
VALUES
    ('growth', 'Growth', 'Up to 500 seats, metered usage billed monthly.', 35000, 'monthly', 500),
    ('enterprise', 'Enterprise', 'Dedicated concurrency budget with premium support.', 125000, 'monthly', 5000)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price_cents = EXCLUDED.price_cents,
    billing_period = EXCLUDED.billing_period,
    max_seats = EXCLUDED.max_seats;
