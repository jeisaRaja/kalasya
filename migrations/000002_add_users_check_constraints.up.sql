ALTER TABLE users 
    ADD CONSTRAINT created_at_check CHECK (created_at <= CURRENT_TIMESTAMP);

ALTER TABLE users 
    ADD CONSTRAINT updated_at_check CHECK (updated_at <= CURRENT_TIMESTAMP);
