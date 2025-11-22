ALTER TABLE tasks ADD COLUMN IF NOT EXISTS tenant_id BIGINT REFERENCES tenants(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_tasks_tenant_id ON tasks(tenant_id);

UPDATE tasks SET tenant_id = (
    SELECT tenant_id FROM users WHERE users.id = tasks.created_by_id LIMIT 1
) WHERE tenant_id IS NULL;

ALTER TABLE tasks ALTER COLUMN tenant_id SET NOT NULL;
