CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);