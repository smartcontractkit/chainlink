-- +goose Up

-- Create new user roles enum for users table
CREATE TYPE user_roles AS ENUM ('admin', 'edit', 'run', 'view');

-- Add new role column to users table, type enum
ALTER TABLE users ADD role user_roles NOT NULL DEFAULT 'view';

-- We are migrating up from a single user full access user - this should be reflected as the admin 
UPDATE users SET role = 'admin';

CREATE UNIQUE INDEX unique_users_lowercase_email ON users (lower(email));

-- Update sessions table include email column to key on user tied to session
DELETE FROM sessions;
ALTER TABLE sessions ADD email text NOT NULL;

ALTER TABLE sessions ADD CONSTRAINT sessions_fk_email FOREIGN KEY(email) REFERENCES users(email) ON DELETE cascade;

-- +goose Down

ALTER TABLE users DROP COLUMN role;
DROP TYPE user_roles;

ALTER TABLE sessions DROP CONSTRAINT sessions_fk_email;
ALTER TABLE sessions DROP COLUMN email;

DROP INDEX unique_users_lowercase_email;
