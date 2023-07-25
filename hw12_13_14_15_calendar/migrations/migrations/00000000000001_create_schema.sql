-- +goose Up
CREATE SCHEMA IF NOT EXISTS hw12calendar;
--SET search_path = hw12calendar;
-- +goose Down
--SET search_path = public;
DROP SCHEMA hw12calendar CASCADE;

