CREATE TABLE IF NOT EXISTS pack_sizes (
    id          BIGSERIAL PRIMARY KEY,
    size        INTEGER NOT NULL UNIQUE,
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);


