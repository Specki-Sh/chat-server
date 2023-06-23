CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    room_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'sent',
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (status IN ('sent', 'read')),
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);