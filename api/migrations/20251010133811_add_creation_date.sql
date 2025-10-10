-- +goose Up
ALTER TABLE parts ADD COLUMN creation_date TIMESTAMP;
CREATE INDEX idx_parts_creation_date ON parts(creation_date);

-- +goose Down
DROP INDEX IF EXISTS idx_parts_creation_date;
ALTER TABLE parts DROP COLUMN creation_date;
