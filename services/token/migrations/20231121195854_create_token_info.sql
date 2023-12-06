-- +goose Up
-- +goose StatementBegin
CREATE TABLE token_info (
    id SERIAL PRIMARY KEY,

    ssh TEXT NOT NULL,
    token TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX token_info_token ON token_info(token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX token_info_token;
DROP TABLE token_info;
-- +goose StatementEnd
