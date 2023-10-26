-- +goose Up
CREATE SCHEMA IF NOT EXISTS hw15calendar;
ALTER SCHEMA hw15calendar OWNER TO hw15user;
-- SET search_path = hw15calendar;
-- +goose Down
-- SET search_path = public;
DROP SCHEMA IF EXISTS hw15calendar CASCADE;

