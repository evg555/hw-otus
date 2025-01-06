-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    uuid UUID PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    user_id UUID,
    notification_time INT
);

-- +goose Down
DROP TABLE IF EXISTS events;
