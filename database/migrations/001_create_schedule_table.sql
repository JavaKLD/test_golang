-- +goose Up
CREATE TABLE IF NOT EXISTS schedules (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    aid_name LONGTEXT NOT NULL,
    aid_per_day BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    duration BIGINT,
    created_at DATETIME
);

-- +goose Down
DROP TABLE IF EXISTS schedules;