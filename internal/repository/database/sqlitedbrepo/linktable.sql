CREATE TABLE IF NOT EXISTS `link` (
    `id` INTEGER PRIMARY KEY,
    `user_id` TEXT NOT NULL,
    `short_url` TEXT NOT NULL,
    `original_url` TEXT NOT NULL UNIQUE
);
