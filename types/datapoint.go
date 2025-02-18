package types

type DataPoint struct {
	Timestamp Timestamp `csv:"timestamp" db:"time"      json:"timestamp"`
	FlowRate  float64   `csv:"flow_rate" db:"flow_rate" json:"flowRate"`
}
