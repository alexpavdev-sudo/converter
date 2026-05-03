-- +goose Up
-- +migrate Up
CREATE TABLE files
(
    id             BIGSERIAL PRIMARY KEY,
    stored_name    VARCHAR(128) NOT NULL,
    extension      VARCHAR(20)  NOT NULL,
    original_name  VARCHAR(255) NOT NULL,
    path           TEXT         NOT NULL,
    status         SMALLINT     NOT NULL DEFAULT 0,
    processed_path TEXT,
    format         VARCHAR(50)  NOT NULL,
    size           BIGINT       NOT NULL,
    size_processed BIGINT,
    deleted        SMALLINT     NOT NULL DEFAULT 0,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_files_stored_name ON files (stored_name);
CREATE UNIQUE INDEX idx_files_processed_path ON files (processed_path);

CREATE INDEX idx_files_status ON files (status);

-- +goose Down
DROP TABLE IF EXISTS files;
