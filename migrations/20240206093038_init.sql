-- +goose Up
-- +goose StatementBegin

create or replace function uuid_generate_v7()
    returns uuid
as $$
declare
unix_ts_ms bytea;
    uuid_bytes bytea;
begin
    unix_ts_ms = substring(int8send(floor(extract(epoch from clock_timestamp()) * 1000)::bigint) from 3);

    -- use random v4 uuid as starting point (which has the same variant we need)
    uuid_bytes = uuid_send(gen_random_uuid());

    -- overlay timestamp
    uuid_bytes = overlay(uuid_bytes placing unix_ts_ms from 1 for 6);

    -- set version 7
    uuid_bytes = set_byte(uuid_bytes, 6, (b'0111' || get_byte(uuid_bytes, 6)::bit(4))::bit(8)::int);

return encode(uuid_bytes, 'hex')::uuid;
end
$$
language plpgsql
    volatile;

CREATE TABLE users
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    email       VARCHAR(128) DEFAULT NULL,
    name        VARCHAR(128) DEFAULT NULL,
    birthday    TIMESTAMPTZ(0) DEFAULT NULL,

    created_at  TIMESTAMPTZ(0) DEFAULT NOW(),
    updated_at  TIMESTAMPTZ(0) DEFAULT NOW()
);
CREATE UNIQUE INDEX users_email_idx ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
drop function uuid_generate_v7;
-- +goose StatementEnd
