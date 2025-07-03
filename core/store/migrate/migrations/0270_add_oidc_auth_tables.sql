-- +goose Up
CREATE TABLE IF NOT EXISTS oidc_sessions (
	id text PRIMARY KEY,
    user_email text NOT NULL,
    user_role user_roles,
    created_at timestamp WITH time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS oidc_user_api_tokens (
    user_email text PRIMARY KEY,
    user_role user_roles,
    token_key text UNIQUE NOT NULL,
    token_salt text NOT NULL,
    token_hashed_secret text NOT NULL,
    created_at timestamp WITH time zone NOT NULL
);

-- +goose Down
DROP TABLE oidc_sessions;
DROP TABLE oidc_user_api_tokens;
