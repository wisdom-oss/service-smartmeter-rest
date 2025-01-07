-- name: get-data
SELECT
    *
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1;

-- name: get-data-from
SELECT
    *
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1
    AND "time">$2;

-- name: get-data-until
SELECT
    *
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1
    AND "time"<$2;

-- name: get-data-in-range
SELECT
    *
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1
    AND "time">$2
    AND "time"<$3;