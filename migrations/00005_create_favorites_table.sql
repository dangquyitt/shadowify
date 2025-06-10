-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS favorites (
        id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id TEXT NOT NULL,
        video_id TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (user_id, video_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS favorites;

-- +goose StatementEnd