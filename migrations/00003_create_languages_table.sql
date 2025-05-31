-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS languages (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        code TEXT UNIQUE NOT NULL DEFAULT '',
        flag_url TEXT NOT NULL DEFAULT '',
        name TEXT NOT NULL DEFAULT '',
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

ALTER TABLE videos
ADD COLUMN IF NOT EXISTS language_id uuid;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS languages;

-- +goose StatementEnd