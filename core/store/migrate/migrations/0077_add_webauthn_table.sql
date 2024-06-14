-- +goose Up
CREATE TABLE web_authns (
    "id" BIGSERIAL PRIMARY KEY,
    "email" text NOT NULL,
    "public_key_data" jsonb NOT NULL,
    CONSTRAINT fk_email
        FOREIGN KEY(email)
        REFERENCES users(email)
);

CREATE UNIQUE INDEX web_authns_email_idx ON web_authns (lower(email));

-- +goose Down
DROP TABLE IF EXISTS web_authns;
