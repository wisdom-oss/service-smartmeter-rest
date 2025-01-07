package types

import (
	"encoding/json"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var geoJsonMarshalOptions = []geojson.EncodeGeometryOption{
	geojson.EncodeGeometryWithBBox(),
	geojson.EncodeGeometryWithMaxDecimalDigits(15), //nolint:mnd
}

// SmartMeter represents real world smart meter.
type SmartMeter struct {
	ID                   string                 `db:"id"                    json:"id"`
	Geometry             geom.T                 `db:"geometry"`
	Name                 *string                `db:"name"                  json:"name"`
	Key                  string                 `db:"key"                   json:"key"`
	AdditionalProperties map[string]interface{} `db:"additional_properties" json:"additionalProperties"`
}

// MarshalJSON implements the [json.Marshaler] interface and is used to convert
// the database depiction of the geometry to GeoJSON.
func (s SmartMeter) MarshalJSON() ([]byte, error) {
	var out struct {
		SmartMeter
		Geometry json.RawMessage `json:"geometry"`
	}
	encodedGeometry, err := geojson.Marshal(s.Geometry, geoJsonMarshalOptions...)
	if err != nil {
		return nil, err
	}
	out.ID = s.ID
	out.Name = s.Name
	out.Key = s.Key
	out.AdditionalProperties = s.AdditionalProperties
	out.Geometry = encodedGeometry
	return json.Marshal(out)
}
