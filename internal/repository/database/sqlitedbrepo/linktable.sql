CREATE TABLE IF NOT EXISTS `links` (
    `id` INTEGER PRIMARY KEY,
    `user_id` TEXT NOT NULL,
    `short_url` TEXT NOT NULL,
    `original_url` TEXT NOT NULL UNIQUE,
    `is_deleted` INTEGER NOT NULL DEFAULT 0
);
