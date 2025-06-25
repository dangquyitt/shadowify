-- +goose Up
-- +goose StatementBegin
ALTER TABLE sentences
ADD COLUMN IF NOT EXISTS meaning_en TEXT NOT NULL DEFAULT '';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE sentences
DROP COLUMN IF EXISTS meaning_en;

-- +goose StatementEnd