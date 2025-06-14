-- +goose Up
-- +goose StatementBegin
ALTER TABLE videos
ADD COLUMN IF NOT EXISTS view_count BIGINT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE videos
DROP COLUMN IF EXISTS view_count;

-- +goose StatementEnd