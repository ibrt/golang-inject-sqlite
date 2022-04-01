-- +migrate Up

CREATE TABLE second_table
(
    id    TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);

-- +migrate Down

DROP TABLE second_table;