CREATE TABLE IF NOT EXISTS tasks (
    id          BIGSERIAL    PRIMARY KEY,
    title       TEXT         NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    status      TEXT         NOT NULL,
    recurrence  JSONB,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_status     ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks (created_at DESC);

-- Индекс по типу повторения внутри JSONB
CREATE INDEX IF NOT EXISTS idx_tasks_recurrence_type
    ON tasks ((recurrence->>'type'))
    WHERE recurrence IS NOT NULL;
