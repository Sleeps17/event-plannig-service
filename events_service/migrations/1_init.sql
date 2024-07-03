-- +goose Up
CREATE SCHEMA IF NOT EXISTS events_schema;

CREATE TABLE IF NOT EXISTS events_schema.events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    room_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS events_schema.participants (
    event_id INTEGER NOT NULL,
    employee_id INTEGER NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events_schema.events (id),
    PRIMARY KEY (event_id, employee_id)
);

-- +goose Down
DROP TABLE IF EXISTS events_schema.events;
DROP TABLE IF EXISTS events_schema.participants;
DROP SCHEMA IF EXISTS events_schema;