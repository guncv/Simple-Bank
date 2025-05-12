ALTER TABLE users
ADD COLUMN role VARCHAR NOT NULL DEFAULT 'depositor';

ALTER TABLE users
ADD CONSTRAINT check_role CHECK (role IN ('depositor', 'banker'));
