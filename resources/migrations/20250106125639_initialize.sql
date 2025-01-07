-- +goose Up
-- +goose StatementBegin
-- activate timescaledb if not already enabled
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- activate postgis if not already enabled
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE EXTENSION IF NOT EXISTS postgis_sfcgal;

-- create the geodata schema and required tables if it doesn't exists already
CREATE SCHEMA IF NOT EXISTS geodata;

CREATE TABLE IF NOT EXISTS
    geodata.layers (
        id UUID DEFAULT gen_random_uuid () NOT NULL PRIMARY KEY,
        NAME TEXT NOT NULL,
        description TEXT,
        "table" TEXT NOT NULL UNIQUE,
        crs INT NOT NULL,
        attribution TEXT,
        private BOOLEAN DEFAULT FALSE
    );

CREATE TABLE IF NOT EXISTS
    geodata.smartmeters (
        id bigserial PRIMARY KEY,
        geometry geometry NOT NULL,
        "key" TEXT NOT NULL,
        "name" TEXT,
        additional_properties jsonb DEFAULT NULL
    );

INSERT INTO
    geodata.layers (NAME, "table", crs)
VALUES
    ('Smartmeter Locations', 'smartmeters', 4326);

-- now setup the timeseries schema
CREATE SCHEMA IF NOT EXISTS timeseries;

-- create the timeseries table
CREATE TABLE IF NOT EXISTS
    timeseries.smartmeter_data (
        "time" timestamptz NOT NULL DEFAULT NOW(),
        smart_meter INT NOT NULL REFERENCES geodata.smartmeters (id) ON UPDATE CASCADE ON DELETE CASCADE,
        flow_rate DOUBLE PRECISION DEFAULT 0
    );

-- convert it into a hypertable
SELECT
    create_hypertable ('timeseries.smartmeter_data', by_range ('time', INTERVAL '7 day'));

-- configure the auth schema as we need to register ourselves in the service
-- database
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS
    auth.services (
        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
        "name" TEXT NOT NULL UNIQUE,
        description TEXT,
        supported_scope_levels scope_level[]
    );

INSERT INTO
    auth.services ("name", supported_scope_levels)
VALUES
    ('smartmeters', '{read, write, *}') ON CONFLICT
DO NOTHING
-- +goose StatementEnd