-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders(
    id serial PRIMARY KEY,
    data jsonb NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
