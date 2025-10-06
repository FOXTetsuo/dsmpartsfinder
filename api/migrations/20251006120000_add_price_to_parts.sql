-- +goose Up
-- +goose StatementBegin
ALTER TABLE parts ADD COLUMN price TEXT;
UPDATE parts SET price = '' WHERE price IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE parts DROP COLUMN price;
-- +goose StatementEnd
