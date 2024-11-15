ALTER TABLE users
    DROP CONSTRAINT IF EXISTS created_at_check;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS updated_at_check;
