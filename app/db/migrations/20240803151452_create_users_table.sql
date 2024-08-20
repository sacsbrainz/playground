-- +goose Up
create table if not exists users (
    id uuid primary key,
    first_name text not null,
    last_name text not null,
    email text unique not null,
    password_hash text not null,
    avatar text default null,
    org_id integer references organizations default null,
    last_login datetime default null,
    failed_login_attempt int default 0,
    email_verified_at datetime default null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp,
    deleted_at datetime default null
);
-- +goose Down
drop table if exists users;