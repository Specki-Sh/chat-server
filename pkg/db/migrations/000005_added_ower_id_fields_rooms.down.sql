ALTER TABLE rooms
    DROP CONSTRAINT rooms_owner_id_fkey,
    DROP COLUMN owner_id;