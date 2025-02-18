package types

import (
	"encoding/json"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

const coordinatePrecision = 15

var geoJsonMarshalOptions = []geojson.EncodeGeometryOption{
	geojson.EncodeGeometryWithBBox(),
	geojson.EncodeGeometryWithMaxDecimalDigits(coordinatePrecision),
}

// Smartmeter represents real world smart meter.
type Smartmeter struct {
	ID                   int            `db:"id"`
	Geometry             geom.T         `db:"geometry"`
	Name                 *string        `db:"name"`
	Key                  string         `db:"key"`
	AdditionalProperties map[string]any `db:"additional_properties"`
}

// MarshalJSON implements the [json.Marshaler] interface and is used to convert
// the database depiction of the geometry to GeoJSON.
func (s Smartmeter) MarshalJSON() ([]byte, error) {
	var out struct {
		ID                   int             `json:"id"`
		Geometry             json.RawMessage `json:"geometry"`
		Name                 *string         `json:"name"`
		Key                  string          `json:"key"`
		AdditionalProperties map[string]any  `json:"additionalProperties"`
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

func (s *Smartmeter) UnmarshalJSON(src []byte) error {
	var smartmeter Smartmeter
	var in struct {
		ID                   int             `json:"id"`
		Geometry             json.RawMessage `json:"geometry"`
		Name                 *string         `json:"name"`
		Key                  string          `json:"key"`
		AdditionalProperties map[string]any  `json:"additionalProperties"`
	}
	if err := json.Unmarshal(src, &in); err != nil {
		return err
	}

	var decodedGeometry geom.T
	if err := geojson.Unmarshal(in.Geometry, &decodedGeometry); err != nil {
		return err
	}

	smartmeter.ID = in.ID
	smartmeter.Name = in.Name
	smartmeter.Key = in.Key
	smartmeter.Geometry = decodedGeometry
	smartmeter.AdditionalProperties = in.AdditionalProperties

	*s = smartmeter
	return nil
}
