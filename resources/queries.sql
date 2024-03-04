-- name: timeseries-overview
SELECT smart_meter,
    MAX(time) AS last_entry,
    MIN(time) AS first_entry
FROM time_series.smart_meter_data
GROUP BY smart_meter;

-- name: timeseries-exists
SELECT EXISTS(SELECT DISTINCT smart_meter_data.smart_meter FROM time_series.smart_meter_data WHERE smart_meter = $1);

-- name: timeseries
SELECT time, flow_rate
FROM time_series.smart_meter_data
WHERE smart_meter = $1;

-- name: timeseries-daterange
SELECT time, flow_rate
FROM time_series.smart_meter_data
WHERE smart_meter = $1
AND time > $2
AND time < $3;

-- name: timeseries-daterange-from
SELECT time, flow_rate
FROM time_series.smart_meter_data
WHERE smart_meter = $1
  AND time > $2
ORDER BY time;

-- name: timeseries-daterange-until
SELECT time, flow_rate
FROM time_series.smart_meter_data
WHERE smart_meter = $1
  AND time < $2
ORDER BY time;