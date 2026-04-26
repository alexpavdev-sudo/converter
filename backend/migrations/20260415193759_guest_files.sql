-- +goose Up
CREATE TABLE guest_files
(
    guest_id BIGINT NOT NULL,
    file_id  BIGINT NOT NULL,
    PRIMARY KEY (guest_id, file_id),
    CONSTRAINT fk_guest_files_guest
        FOREIGN KEY (guest_id)
            REFERENCES guests (id)
            ON DELETE CASCADE
            ON UPDATE RESTRICT,
    CONSTRAINT fk_guest_files_file
        FOREIGN KEY (file_id)
            REFERENCES files (id)
            ON DELETE CASCADE
            ON UPDATE RESTRICT
);

-- Индексы для быстрого поиска
CREATE INDEX idx_guest_files_guest_id ON guest_files (guest_id);
CREATE INDEX idx_guest_files_file_id ON guest_files (file_id);

-- +goose Down
DROP TABLE IF EXISTS guest_files;
