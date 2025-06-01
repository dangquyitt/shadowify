-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS videos (
        id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
        title TEXT NOT NULL DEFAULT '',
        full_title TEXT NOT NULL DEFAULT '',
        description TEXT NOT NULL DEFAULT '',
        thumbnail TEXT NOT NULL DEFAULT '',
        duration BIGINT NOT NULL DEFAULT 0,
        duration_string TEXT NOT NULL DEFAULT '',
        youtube_id TEXT UNIQUE NOT NULL DEFAULT '',
        tags jsonb NOT NULL DEFAULT '{}'::jsonb,
        categories jsonb NOT NULL DEFAULT '{}'::jsonb,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS videos;

-- +goose StatementEnd