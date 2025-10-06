-- +goose Up
-- +goose StatementBegin
INSERT INTO sites (id, site_name, site_url)
VALUES (1, 'SchadeAutos', 'https://www.schadeautos.nl')
ON CONFLICT(id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sites WHERE id = 1;
-- +goose StatementEnd
