-- +goose Up
CREATE TABLE IF NOT EXISTS ldap_sessions (
	id text PRIMARY KEY,	
    user_email text NOT NULL,
    user_role user_roles,
    localauth_user BOOLEAN,
    created_at timestamp with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS ldap_user_api_tokens (
    user_email text PRIMARY KEY,
    user_role user_roles,
    localauth_user BOOLEAN,
    token_key text UNIQUE NOT NULL,
    token_salt text NOT NULL,
    token_hashed_secret text NOT NULL,
    created_at timestamp with time zone NOT NULL
);

-- +goose Down
DROP TABLE ldap_sessions;
DROP TABLE ldap_user_api_tokens;
