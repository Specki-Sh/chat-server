CREATE TABLE members (
    room_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);