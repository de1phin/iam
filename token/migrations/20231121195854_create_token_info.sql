-- +goose Up
-- +goose StatementBegin
CREATE TABLE token_info (
    id serial primary key,

    account_id varchar(512) NOT NULL,
    token UUID NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT (NOW() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() at time zone 'utc'),
    deleted_at TIMESTAMP,

    created_name varchar(256),
    updated_name varchar(256),
    deleted_name varchar(256)
)

CREATE UNIQUE INDEX token_info_token_username ON token_info(token, account_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX token_info_token_username;
DROP TABLE token_info;
-- +goose StatementEnd
