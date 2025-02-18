package meters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice/types"
)

func testCreation(t *testing.T) {
	path := "" // path is empty here due to the format string already introducing the slash

	for _, testMeter := range testMeters {
		body, err := json.Marshal(testMeter)
		assert.NoError(t, err)

		r := bytes.NewReader(body)

		req := httptest.NewRequest("PUT", fmt.Sprintf("%s/%s", routePrefix, path), r)
		res := httptest.NewRecorder()

		router.Handler().ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)

		var meter types.Smartmeter
		err = json.NewDecoder(res.Result().Body).Decode(&meter)
		assert.NoError(t, err)

		assert.Equal(t, testMeter.Name, *meter.Name)
		assert.Equal(t, testMeter.Longitude, meter.Geometry.FlatCoords()[0])
		assert.Equal(t, testMeter.Latitude, meter.Geometry.FlatCoords()[1])
		assert.Equal(t, testMeter.AdditionalParameters, meter.AdditionalProperties)
	}
}
