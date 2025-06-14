-- +goose Up
-- +goose StatementBegin
ALTER TABLE segments
ADD COLUMN IF NOT EXISTS cefr TEXT NOT NULL DEFAULT '';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE segments
DROP COLUMN IF EXISTS cefr;

-- +goose StatementEnd