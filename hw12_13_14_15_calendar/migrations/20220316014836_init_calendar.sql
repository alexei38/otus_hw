-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
     id serial PRIMARY KEY,
     title TEXT NOT NULL,
     start timestamptz NOT NULL,
     stop timestamptz NOT NULL,
     description TEXT,
     user_id int NOT NULL,
     before_send bigint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd
