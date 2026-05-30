-- =====================================================================
-- pdf-reader backend :: PostgreSQL schema initialization
-- Derived from documentation/ERD.jpg
-- =====================================================================

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE ai_job_type_enum AS ENUM ('TRANSLATE', 'SUMMARIZE', 'COMMENT');
CREATE TYPE ai_job_status_enum AS ENUM ('pending', 'processing', 'completed', 'failed');

CREATE TABLE roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(64)  NOT NULL UNIQUE,
    code        VARCHAR(64)  NOT NULL UNIQUE,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);

CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id     UUID         NOT NULL REFERENCES roles(id),
    fullname    VARCHAR(128),
    username    VARCHAR(64)  NOT NULL UNIQUE,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);

CREATE TABLE workspaces (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id),
    name        VARCHAR(255) NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);

CREATE TABLE files (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID         NOT NULL REFERENCES users(id),
    workspace_id        UUID         NOT NULL REFERENCES workspaces(id),
    original_file_name  VARCHAR(512) NOT NULL,
    name                VARCHAR(512) NOT NULL,
    url                 VARCHAR(1024) NOT NULL,
    type                VARCHAR(32)  NOT NULL DEFAULT 'pdf',

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);

CREATE TABLE comments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id     UUID NOT NULL REFERENCES files(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    message     TEXT NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);

CREATE TABLE ai_jobs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id),
    workspace_id  UUID NOT NULL REFERENCES workspaces(id),
    file_id       UUID NOT NULL REFERENCES files(id),
    type          ai_job_type_enum   NOT NULL,
    status        ai_job_status_enum NOT NULL DEFAULT 'pending',
    input         TEXT,
    output        TEXT,
    duration      BIGINT NOT NULL DEFAULT 0,
    retry_count   INT    NOT NULL DEFAULT 0,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by  UUID,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by  UUID,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMPTZ,
    deleted_by  UUID
);

CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_workspaces_user_id ON workspaces(user_id);
CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_workspace_id ON files(workspace_id);
CREATE INDEX idx_comments_file_id ON comments(file_id);
CREATE INDEX idx_ai_jobs_user_id ON ai_jobs(user_id);
CREATE INDEX idx_ai_jobs_workspace_id ON ai_jobs(workspace_id);
CREATE INDEX idx_ai_jobs_file_id ON ai_jobs(file_id);

COMMIT;
