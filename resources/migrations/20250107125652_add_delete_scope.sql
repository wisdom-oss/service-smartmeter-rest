-- +goose Up
-- +goose StatementBegin
UPDATE auth.services
SET
    supported_scope_levels='{read, write, delete, *}'
WHERE
    "name"='smartmeters';

-- +goose StatementEnd