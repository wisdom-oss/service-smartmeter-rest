package types

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jszwec/csvutil"
)

// invalidTimestamp is a predefined value for marking invalid timestamp formats
// coming in while parsing csv.
var invalidTimestamp = Timestamp{Timestamptz: pgtype.Timestamptz{Valid: false}}

// allowedTimestampFormats contains all allowed and parseable timestamp foramts
// for reading in a timestamp.
var allowedTimestampFormats = []string{time.RFC3339Nano, time.RFC3339, `02.01.2006 15:04`}

// Timestamp embeds a [pgtype.Timestamptz] to allow reading database entries
// pgx handling.
type Timestamp struct{ pgtype.Timestamptz }

// MarshalCSV implements the [csvutil.Marshaler] interface which allows writing
// [pgtype.Timestamptz] values into a csv file.
func (t Timestamp) MarshalCSV() ([]byte, error) {
	return csvutil.Marshal(t.Time)
}

// Unmarshal implements the [csvutils.Unmarshaler] interface which enables the
// parsing of timestamps for incoming CSV files and bodies.
func (t *Timestamp) UnmarshalCSV(data []byte) error {
	if len(data) == 0 {
		*t = invalidTimestamp
		return nil
	}

	var timestamp time.Time
	var err error
	for _, allowedTimestampFormat := range allowedTimestampFormats {
		timestamp, err = time.Parse(allowedTimestampFormat, string(data))
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	return t.Scan(timestamp)
}
