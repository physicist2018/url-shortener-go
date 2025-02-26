CREATE TABLE IF NOT EXISTS links (
    user_id VARCHAR(36) NOT NULL,
    short_url VARCHAR(36) NOT NULL,
    original_url VARCHAR(512) NOT NULL UNIQUE,
    is_deleted BOOLEAN DEFAULT FALSE
);
