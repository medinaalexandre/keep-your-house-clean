ALTER TABLE compliments ADD COLUMN IF NOT EXISTS viewed_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_compliments_viewed_at ON compliments(viewed_at);

