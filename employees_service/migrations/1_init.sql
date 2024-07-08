-- +goose Up
CREATE SCHEMA IF NOT EXISTS employees_schema;

CREATE TABLE IF NOT EXISTS employees_schema.employees(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(20),
    last_name VARCHAR(20),
    patronymic VARCHAR(20),
    email VARCHAR(50),
    mobile_number VARCHAR(10)
);

-- +goose Down
DROP TABLE IF EXISTS employees_schema.employees;
DROP SCHEMA IF EXISTS employees_schema;