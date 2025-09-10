CREATE TABLE IF NOT EXISTS user_applications
(
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    text       TEXT NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'new',
    file_url   VARCHAR(120),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
    );
