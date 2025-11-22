CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    scheduled_to TIMESTAMP,
    scheduled_by_id BIGINT,
    frequency_value INTEGER NOT NULL,
    frequency_unit VARCHAR(20) NOT NULL CHECK (frequency_unit IN ('days', 'weeks', 'months')),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_by_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by_id BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by_id BIGINT,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
CREATE INDEX IF NOT EXISTS idx_tasks_created_by_id ON tasks(created_by_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed);
