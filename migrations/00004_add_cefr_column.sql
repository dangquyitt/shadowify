-- +goose Up
-- +goose StatementBegin
ALTER TABLE videos
ADD COLUMN IF NOT EXISTS cefr TEXT NOT NULL DEFAULT '';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE videos
DROP COLUMN IF EXISTS cefr;

-- +goose StatementEnd