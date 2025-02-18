package meters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice/types"
)

func listMeters(t *testing.T) {
	path := "" // path is empty here due to the format string already introducing the slash

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/%s", routePrefix, path), nil)
	res := httptest.NewRecorder()

	router.Handler().ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode) // assert that the request has been successful
	if t.Failed() {
		return
	}

	var readMeters []types.Smartmeter
	err := json.NewDecoder(res.Body).Decode(&readMeters)
	assert.NoError(t, err)

	if t.Failed() {
		return
	}

	for _, readMeter := range readMeters {
		testMeterIdx := slices.IndexFunc(testMeters, func(testMeterParams newMeterParameters) bool {
			return *readMeter.Name == testMeterParams.Name
		})
		assert.NotEqual(t, -1, testMeterIdx) // nolint:mnd
		if t.Failed() {
			return
		}
		testMeter := testMeters[testMeterIdx]

		assert.Equal(t, testMeter.Name, *readMeter.Name)
		assert.Equal(t, testMeter.Longitude, readMeter.Geometry.FlatCoords()[0])
		assert.Equal(t, testMeter.Latitude, readMeter.Geometry.FlatCoords()[1])
		assert.Equal(t, testMeter.AdditionalParameters, readMeter.AdditionalProperties)
	}

}
