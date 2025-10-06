-- +goose Up
-- +goose StatementBegin
ALTER TABLE parts ADD COLUMN last_seen DATETIME;
UPDATE parts SET last_seen = created_at WHERE last_seen IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE parts DROP COLUMN last_seen;
-- +goose StatementEnd
