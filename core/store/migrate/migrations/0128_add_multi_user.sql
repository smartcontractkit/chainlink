-- +goose Up

-- Create new user roles enum for users table
CREATE TYPE user_roles AS ENUM ('admin', 'edit', 'edit_minimal', 'view');

-- Add new role column to users table, type enum
ALTER TABLE users ADD role user_roles NOT NULL DEFAULT 'view';

-- Update sessions table include email column to key on user tied to session
ALTER TABLE sessions ADD email text NOT NULL;
ALTER TABLE sessions ADD CONSTRAINT sessions_fk_email FOREIGN KEY(email) REFERENCES users(email);

-- +goose Down

ALTER TABLE users DROP COLUMN role;
DROP TYPE user_roles;

ALTER TABLE sessions DROP COLUMN email;
