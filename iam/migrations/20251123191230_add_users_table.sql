-- +goose Up
create table users (
   id uuid primary key default gen_random_uuid(),
    login varchar not null unique,
    email varchar not null unique,
	notification_method jsonb not null default '{}',
    password_hash varchar not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +goose Down
drop table users;

