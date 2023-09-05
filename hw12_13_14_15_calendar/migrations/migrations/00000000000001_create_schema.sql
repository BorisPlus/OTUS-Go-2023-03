-- +goose Up
CREATE SCHEMA IF NOT EXISTS hw15calendar;
-- SET search_path = hw15calendar;
-- +goose Down
-- SET search_path = public;
DROP SCHEMA hw15calendar CASCADE;

