CREATE TABLE IF NOT EXISTS couriers (
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT      NOT NULL,
    loc_lat FLOAT8    NOT NULL,
    loc_lon FLOAT8    NOT NULL,
    status  TEXT      NOT NULL DEFAULT 'free'
);

CREATE TABLE IF NOT EXISTS orders (
    id         BIGSERIAL PRIMARY KEY,
    address    TEXT      NOT NULL,
    loc_lat    FLOAT8    NOT NULL,
    loc_lon    FLOAT8    NOT NULL,
    status     TEXT      NOT NULL DEFAULT 'pending',
    courier_id BIGINT    REFERENCES couriers(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_status     ON orders  (status);
CREATE INDEX IF NOT EXISTS idx_couriers_status   ON couriers (status);
