package types

import "github.com/jackc/pgx/v5/pgtype"

type DataPoint struct {
	Timestamp pgtype.Timestamptz `json:"timestamp" db:"time" cbor:"1"`
	Value     pgtype.Float8      `json:"value" db:"flow_rate" cbor:"2"`
}
