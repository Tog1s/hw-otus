-- +goose Up
-- +goose StatementBegin
create extension if not exists "uuid-ossp";
CREATE TABLE if not exists events (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title text NOT null,
    date_time timestamptz NOT null,
    end_time timestamptz NOT null,
    description text NOT null,
    user_id integer NOT null,
    notify_before timestamptz NOT null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
