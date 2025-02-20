CREATE TABLE IF NOT EXISTS `link` (
    `id` SERIAL PRIMARY KEY,
    `user_id` VARCHAR(36) NOT NULL,
    `short_url` VARCHAR(36) NOT NULL,
    `original_url` VARCHAR(512) NOT NULL UNIQUE
);
