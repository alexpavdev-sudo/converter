-- +goose Up
CREATE TABLE errors
(
    id         BIGSERIAL PRIMARY KEY,
    file_id    BIGINT    NOT NULL,
    details    TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_errors_file
        FOREIGN KEY (file_id)
            REFERENCES files (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_errors_file_id ON errors(file_id);
CREATE INDEX IF NOT EXISTS idx_errors_created_at ON errors(created_at);

-- +goose Down
DROP TABLE IF EXISTS errors;