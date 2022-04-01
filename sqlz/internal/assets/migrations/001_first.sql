-- +migrate Up

CREATE TABLE first_table
(
    id    TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);

-- +migrate Down

DROP TABLE first_table;