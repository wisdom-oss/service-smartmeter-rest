-- name: cleanup-smartmeters
DELETE FROM geodata.smartmeters 
WHERE key = $1;