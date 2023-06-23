ALTER TABLE rooms
    ADD COLUMN owner_id INTEGER,
    ADD FOREIGN KEY (owner_id) REFERENCES users(id);