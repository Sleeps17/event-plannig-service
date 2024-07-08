-- +goose Up
CREATE SCHEMA IF NOT EXISTS users_schema;

CREATE SCHEMA IF NOT EXISTS apps_schema;

CREATE TABLE IF NOT EXISTS users_schema.users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users_schema.users (email);

CREATE TABLE IF NOT EXISTS apps_schema.apps (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS admins (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users_schema.users (id),
    app_id INTEGER REFERENCES apps_schema.apps (id)
);

-- +goose Down
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS app;
DROP TABLE IF EXISTS users;
DROP SCHEMA IF EXISTS users_schema;
DROP SCHEMA IF EXISTS apps_schema;
