DROP TABLE IF EXISTS messages;
CREATE TABLE messages (
    id VARCHAR(32) PRIMARY KEY,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);