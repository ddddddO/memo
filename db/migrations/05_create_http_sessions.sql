-- +migrate Up
CREATE TABLE http_sessions (
    id BIGSERIAL PRIMARY KEY,
    key BYTEA,
    data BYTEA,
    created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_on TIMESTAMPTZ,
    expires_on TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE http_sessions;
