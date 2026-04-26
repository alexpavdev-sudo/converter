-- +goose Up
CREATE TABLE guests
(
    id           BIGSERIAL PRIMARY KEY,
    personal_dir VARCHAR(128) NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_guests_personal_dir ON guests (personal_dir);

-- +goose Down
DROP TABLE IF EXISTS guests;
