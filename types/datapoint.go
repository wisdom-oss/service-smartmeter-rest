package types

import "github.com/jackc/pgx/v5/pgtype"

type DataPoint struct {
	SmartMeter int                `db:"smart_meter" json:"smartMeter"`
	Timestamp  pgtype.Timestamptz `db:"time"        json:"timestamp"`
	FlowRate   float64            `db:"flow_rate"   json:"flowRate"`
}
