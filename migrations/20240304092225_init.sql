-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY, url TEXT);
CREATE INDEX IF NOT EXISTS idx_url ON links(url);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS links;

