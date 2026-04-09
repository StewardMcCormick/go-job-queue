-- +goose Up
-- +goose StatementBegin
CREATE TYPE task_status AS ENUM (
    'TASK_STATUS_COMPLETED',
    'TASK_STATUS_CANCELLED'
    );

CREATE TYPE task_priority AS ENUM (
    'TASK_PRIORITY_BACKGROUND',
    'TASK_PRIORITY_NORMAL',
    'TASK_PRIORITY_HIGH',
    'TASK_PRIORITY_IMMEDIATE'
    );

CREATE TYPE worker_status AS ENUM (
    'WORKER_STATUS_REGISTERED',
    'WORKER_STATUS_WORKING',
    'WORKER_STATUS_DEAD',
    'WORKER_STATUS_UNREGISTERED'
    );

CREATE TABLE tasks
(
    id                  VARCHAR(32) PRIMARY KEY,
    status              task_status   NOT NULL,
    priority            task_priority NOT NULL,
    task_type           VARCHAR(100)  NOT NULL,
    payload             BYTEA,
    should_retry_number INT           NOT NULL,
    retries             INT           NOT NULL,
    deadline            TIMESTAMPTZ,
    created_at          TIMESTAMPTZ   NOT NULL,
    started_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    depends_on          TEXT[]        NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_tasks_completed_at ON tasks (completed_at DESC);
CREATE INDEX idx_tasks_type_completed ON tasks (task_type, completed_at DESC);

CREATE TABLE workers
(
    id              VARCHAR(32) PRIMARY KEY,
    addr            VARCHAR(255)  NOT NULL,
    task_type       VARCHAR(100)  NOT NULL,
    concurrency     INT           NOT NULL,
    status          worker_status NOT NULL,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    unregistered_at TIMESTAMPTZ
);

CREATE INDEX idx_workers_status ON workers (status);
CREATE INDEX idx_workers_task_type ON workers (task_type);

CREATE TABLE task_worker_history
(
    id                    BIGSERIAL PRIMARY KEY,
    task_id               VARCHAR(36) NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    worker_id             VARCHAR(36) NOT NULL REFERENCES workers (id) ON DELETE CASCADE,
    assigned_at           TIMESTAMPTZ NOT NULL,
    completed_at          TIMESTAMPTZ,
    execution_duration_ms INT
);

CREATE INDEX idx_task_worker_history_worker ON task_worker_history (worker_id, completed_at DESC);
CREATE INDEX idx_task_worker_history_task ON task_worker_history (task_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS task_worker_history;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS workers;

DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS task_priority;
DROP TYPE IF EXISTS worker_status;

-- +goose StatementEnd