package types

import "strings"

// OutputFormat represents the format of the output data.
//
// The FromString method allows converting a string to an OutputFormat value.
//
// Usage example:
//
//	var outputFormat OutputFormat
//	outputFormat.FromString("json")
type OutputFormat uint

const (
	JSON OutputFormat = iota
	CSV
	CBOR
)

// FromString takes a string and sets the value of the OutputFormat
// receiver according to the string value. If the string is one of "json",
// "csv", or "cbor", the corresponding OutputFormat value will be assigned.
// Otherwise, the OutputFormat will be set to the default value JSON.
func (f *OutputFormat) FromString(s string) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "json":
		*f = JSON
	case "csv":
		*f = CSV
	case "cbor":
		*f = CBOR
	default:
		*f = JSON
	}
}
