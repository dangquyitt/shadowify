-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

CREATE TABLE
    videos (id)
    -- +goose StatementEnd
    -- +goose Down
    -- +goose StatementBegin
SELECT
    'down SQL query';

-- +goose StatementEnd