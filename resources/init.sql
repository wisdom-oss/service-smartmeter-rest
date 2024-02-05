/*
 * This SQL script contains all commands which prepare the database schema and
 * other tables that are needed to run this microservice
 */

-- name: 01
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- name: 02
CREATE SCHEMA IF NOT EXISTS time_series;

-- name: 03
CREATE TABLE IF NOT EXISTS time_series.smart_meter_data
(
    time        timestamptz      NOT NULL,
    -- smart_meter contains the unique identifier of the smart meter
    smart_meter text             NOT NULL,
    -- flow_rate contains the recorded flow rate in cubic meters (m³)
    flow_rate   double precision NOT NULL
);

-- name: 04
COMMENT ON COLUMN time_series.smart_meter_data.smart_meter IS 'contains the unique identifier of the smart meter';

-- name: 05
COMMENT ON COLUMN time_series.smart_meter_data.flow_rate IS 'contains the recorded flow rate in cubic meters (m³)';

-- name: 06
SELECT create_hypertable('time_series.smart_meter_data', by_range('time'));