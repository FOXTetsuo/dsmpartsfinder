-- +goose Up
-- +goose StatementBegin
INSERT INTO sites (id, site_name, site_url)
VALUES (3, 'Ebay', 'https://www.ebay.nl')
ON CONFLICT(id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sites WHERE id = 3;
-- +goose StatementEnd
