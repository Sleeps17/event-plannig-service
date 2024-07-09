-- +goose Up
CREATE SCHEMA IF NOT EXISTS rooms_schema;

CREATE TABLE IF NOT EXISTS rooms_schema.rooms (
    id SERIAL PRIMARY KEY,
    room_name TEXT NOT NULL,
    capacity INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS rooms_schema.rooms;
DROP SCHEMA IF EXISTS rooms_schema;