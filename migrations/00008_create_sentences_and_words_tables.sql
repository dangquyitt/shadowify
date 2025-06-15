-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS sentences (
        id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id TEXT NOT NULL,
        segment_id TEXT NOT NULL,
        meaning_vi TEXT NOT NULL DEFAULT '',
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (user_id, segment_id)
    );

CREATE TABLE
    IF NOT EXISTS words (
        id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id TEXT NOT NULL,
        meaning_en TEXT NOT NULL,
        meaning_vi TEXT NOT NULL DEFAULT '',
        segment_id TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        UNIQUE (user_id, meaning_en)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sentences;

DROP TABLE IF EXISTS words;

-- +goose StatementEnd