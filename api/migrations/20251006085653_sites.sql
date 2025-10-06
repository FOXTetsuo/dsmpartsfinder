-- +goose Up
CREATE TABLE sites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_url TEXT NOT NULL UNIQUE,
    site_name TEXT NOT NULL
);

-- +goose Down
DROP TABLE sites;
