-- Не стал усложнять задачу наличием внешних ключей.
--  Up:
--      CREATE TABLE IF NOT EXISTS hw12calendar.owners ( 
--          pk BIGSERIAL PRIMARY KEY,
--          name VARCHAR(50) NOT NULL CHECK (name <> '') UNIQUE,
--          contact VARCHAR(50) NOT NULL CHECK (name <> '') UNIQUE
--      );
--      и
--      ... 
--      owner_pk BIGINT NOT NULL,
--      ... 
--      FOREIGN KEY(owner_pk) 
--      REFERENCES hw12calendar.owners(pk)
--      ON DELETE CASCADE
--      ON UPDATE NO ACTION
--      ... 
--  Down:
--      DROP TABLE IF EXISTS hw12calendar.owners;
-- и использование Postgres-типа данных `interval` имеет сложность перевода в time.Duration 
-- +goose Up
CREATE TABLE IF NOT EXISTS hw15calendar.events ( 
 pk BIGSERIAL PRIMARY KEY,
 title VARCHAR NOT NULL CHECK (title <> ''),
 startAt TIMESTAMP WITH TIME ZONE NOT NULL,
 durationSeconds BIGINT NOT NULL DEFAULT 0,
 description TEXT NOT NULL DEFAULT '',
 notifyEarlySeconds BIGINT,
 "owner" VARCHAR NOT NULL CHECK (title <> ''),
 sheduled BOOLEAN DEFAULT FALSE
);
-- +goose Down
DROP TABLE IF EXISTS hw15calendar.events;
