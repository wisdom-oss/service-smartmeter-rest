package types

import "github.com/jackc/pgx/v5/pgtype"

type TimeseriesInformation struct {
	Smartmeter pgtype.Text        `json:"dataSeriesID" db:"smart_meter"`
	FirstEntry pgtype.Timestamptz `json:"firstEntry"`
	LastEntry  pgtype.Timestamptz `json:"lastEntry"`
}
