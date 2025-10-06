-- +goose Up
CREATE TABLE parts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    part_id TEXT NOT NULL,
    description TEXT NOT NULL,
    type_name TEXT NOT NULL,
    name TEXT NOT NULL,
    image_base64 TEXT,
    url TEXT NOT NULL,
    site_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE,
    UNIQUE(part_id, site_id)
);

CREATE INDEX idx_parts_site_id ON parts(site_id);
CREATE INDEX idx_parts_part_id ON parts(part_id);
CREATE INDEX idx_parts_name ON parts(name);

-- +goose Down
DROP INDEX IF EXISTS idx_parts_name;
DROP INDEX IF EXISTS idx_parts_part_id;
DROP INDEX IF EXISTS idx_parts_site_id;
DROP TABLE parts;
