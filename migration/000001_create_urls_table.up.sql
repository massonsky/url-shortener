CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    click_count BIGINT DEFAULT 0
);