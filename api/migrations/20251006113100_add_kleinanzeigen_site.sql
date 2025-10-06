-- +goose Up
-- +goose StatementBegin
INSERT INTO sites (id, site_name, site_url)
VALUES (2, 'Kleinanzeigen', 'https://www.kleinanzeigen.de')
ON CONFLICT(id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sites WHERE id = 2;
-- +goose StatementEnd
