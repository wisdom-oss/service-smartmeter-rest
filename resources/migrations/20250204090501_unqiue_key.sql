-- +goose Up
-- +goose StatementBegin
ALTER TABLE geodata.smartmeters ADD CONSTRAINT uk_key UNIQUE (key);
-- +goose StatementEnd
