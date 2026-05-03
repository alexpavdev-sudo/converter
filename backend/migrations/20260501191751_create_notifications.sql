-- +goose Up
CREATE TABLE notifications
(
    id         BIGSERIAL PRIMARY KEY,
    detail     TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_send    SMALLINT  NOT NULL DEFAULT 0,
    type       SMALLINT  NOT NULL,
    guest_id   BIGINT    NOT NULL,
    CONSTRAINT fk_notifications_guest
        FOREIGN KEY (guest_id)
            REFERENCES guests (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);

CREATE INDEX idx_notifications_unsent_id
    ON notifications (id) WHERE is_send = 0;

CREATE INDEX idx_notifications_unsent_created_at
    ON notifications (created_at) WHERE is_send = 0;

CREATE INDEX idx_notifications_guest_type
    ON notifications (guest_id, type);

-- +goose Down
DROP TABLE IF EXISTS notifications;
