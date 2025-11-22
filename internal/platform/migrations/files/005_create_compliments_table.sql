CREATE TABLE IF NOT EXISTS compliments (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    points INTEGER NOT NULL DEFAULT 0 CHECK (points >= 0 AND points <= 5),
    from_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    to_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by_id BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    deleted_at TIMESTAMP,
    CHECK (from_user_id != to_user_id)
);

CREATE INDEX IF NOT EXISTS idx_compliments_deleted_at ON compliments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_compliments_tenant_id ON compliments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_compliments_from_user_id ON compliments(from_user_id);
CREATE INDEX IF NOT EXISTS idx_compliments_to_user_id ON compliments(to_user_id);
CREATE INDEX IF NOT EXISTS idx_compliments_created_at ON compliments(created_at);

