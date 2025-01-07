-- name: all-meters
SELECT
    *
FROM
    geodata.smartmeters;

-- name: get-meter
SELECT
    *
FROM
    geodata.smartmeters
WHERE
    id=$1;

-- name: insert-meter
INSERT INTO
    geodata.smartmeters ("name", "key", geometry, additional_properties)
VALUES
    ($1, $2, $3, $4);

-- name: delete-meter
DELETE FROM geodata.smartmeters
WHERE
    id=$1;