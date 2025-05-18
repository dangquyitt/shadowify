-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS segments (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        video_id uuid NOT NULL,
        start_sec REAL NOT NULL DEFAULT 0.0,
        end_sec REAL NOT NULL DEFAULT 0.0,
        content TEXT NOT NULL DEFAULT '',
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS segments;

-- +goose StatementEnd