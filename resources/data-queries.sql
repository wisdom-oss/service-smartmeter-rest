-- name: get-data
SELECT
    "time",
    flow_rate
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1;

-- name: get-data-from
SELECT
    "time",
    flow_rate
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1
    AND "time">=$2
ORDER BY "time";

-- name: get-data-until
SELECT
    "time",
    flow_rate
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1
    AND "time"<=$2
ORDER BY "time";

-- name: get-data-in-range
SELECT
    "time",
    flow_rate
FROM
    timeseries.smartmeter_data
WHERE
    smart_meter=$1
    AND "time">=$2
    AND "time"<=$3
ORDER BY "time";