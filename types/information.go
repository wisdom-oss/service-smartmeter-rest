package types

import "github.com/jackc/pgx/v5/pgtype"

type Information struct {
	SmartMeter      int                `db:"snart_meter" json:"smartMeter"`
	FirstRecordTime pgtype.Timestamptz `db:"first_entry" json:"firstRecordTimestamp"`
	LastRecordTime  pgtype.Timestamptz `db:"last_entry"  json:"lastRecordTimestamp"`
}
