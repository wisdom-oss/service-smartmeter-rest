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
    key=$1;

-- name: insert-meter
INSERT INTO
    geodata.smartmeters ("name", "key", geometry, additional_properties)
VALUES
    ($1, $2, ST_Point ($3, $4, 4326), $5)
RETURNING
*;

-- name: delete-meter
DELETE FROM geodata.smartmeters
WHERE
    key=$1;

-- name: meter-exists
SELECT
    EXISTS (
        SELECT
            *
        FROM
            geodata.smartmeters
        WHERE
            key=$1
    );

-- name: update-meter-name
UPDATE geodata.smartmeters
SET
    name=$1
WHERE
    key=$2;

-- name: update-meter-location
UPDATE geodata.smartmeters
SET
    geometry=ST_Point ($1, $2, 4326)
WHERE
    key=$3;

-- name: update-meter-properties
UPDATE geodata.smartmeters
SET
    additional_properties=$1
WHERE
    key=$2;