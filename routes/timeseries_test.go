package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-chi/chi/v5"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
	validator "openapi.tanna.dev/go/validator/openapi3"

	"github.com/wisdom-oss/service-smartmeter-rest/types"
)

func TestTimeSeries_Output_Default(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}

func TestTimeSeries_Output_JSON(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}

func TestTimeSeries_Ouput_CSV(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010?format=csv", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}

func TestTimeSeries_Output_CBOR(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010?format=cbor", nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	// since the validator does not support the cbor content type, this needs to
	// be done manually here
	if responseRecorder.Header().Get("Content-Type") != "application/cbor" {
		t.Errorf("Expected Content-Type to be application/cbor, got %s", responseRecorder.Header().Get("Content-Type"))
	}
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if responseRecorder.Body.Len() == 0 {
		t.Error("Expected non-empty response body")
	}
	var dataseries []types.DataPoint
	err := cbor.NewDecoder(responseRecorder.Body).Decode(&dataseries)
	if err != nil {
		t.Errorf("Error decoding CBOR response body: %v", err)
	}
}

func TestTimeSeries_Error_InvalidFromTimestamp(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010?from=abc", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)
	router.ServeHTTP(responseRecorder, request)

	var body wisdomType.WISdoMError
	err := json.NewDecoder(responseRecorder.Result().Body).Decode(&body)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if !body.Equals(ErrInvalidTimestamp) {
		t.Errorf("Expected error %v, got %v", ErrInvalidTimestamp, body)
	}
}

func TestTimeSeries_ValidFromTimestamp(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	timestamp := "2023-06-01T15:00:00+00:00"
	parsedTime := time.Time{}
	err := parsedTime.UnmarshalText([]byte(timestamp))
	if err != nil {
		t.Errorf("Error scanning timestamp: %v", err)
	}

	request := httptest.NewRequest("GET", "/045010?from="+url.QueryEscape(timestamp), nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)
	router.ServeHTTP(responseRecorder, request)

	var timeseries []types.DataPoint
	err = json.NewDecoder(responseRecorder.Result().Body).Decode(&timeseries)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
		fmt.Fprintf(os.Stdout, "%s\n", responseRecorder.Body)
	}

	if len(timeseries) == 0 {
		t.Errorf("Expected filled timeseries, got %v", timeseries)
	}

	firstPoint := timeseries[0]
	lastPoint := timeseries[len(timeseries)-1]

	if firstPoint.Timestamp.Time.Before(parsedTime) {
		t.Errorf("Expected first data point to be after the provided from timestamp")
	}

	if firstPoint.Timestamp.Time.After(lastPoint.Timestamp.Time) {
		t.Errorf("Expected first data point to be before to the last data point")
	}
}

func TestTimeSeries_ValidUntilTimestamp(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	timestamp := "2023-06-01T15:00:00+00:00"
	parsedTime := time.Time{}
	err := parsedTime.UnmarshalText([]byte(timestamp))
	if err != nil {
		t.Errorf("Error scanning timestamp: %v", err)
	}

	request := httptest.NewRequest("GET", "/045010?until="+url.QueryEscape(timestamp), nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)
	router.ServeHTTP(responseRecorder, request)

	var timeseries []types.DataPoint
	err = json.NewDecoder(responseRecorder.Result().Body).Decode(&timeseries)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
		fmt.Fprintf(os.Stdout, "%s\n", responseRecorder.Body)
	}

	if len(timeseries) == 0 {
		t.Errorf("Expected filled timeseries, got %v", timeseries)
	}

	firstPoint := timeseries[0]
	lastPoint := timeseries[len(timeseries)-1]

	if lastPoint.Timestamp.Time.After(parsedTime) {
		t.Errorf("Expected last data point to be before the provided from timestamp")
	}

	if firstPoint.Timestamp.Time.After(lastPoint.Timestamp.Time) {
		t.Errorf("Expected first data point to be before to the last data point")
	}
}

func TestTimeSeries_ValidTimestamps(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	fromTimestamp := "2023-06-01T15:00:00+00:00"
	parsedFromTime := time.Time{}
	err := parsedFromTime.UnmarshalText([]byte(fromTimestamp))
	if err != nil {
		t.Errorf("Error scanning timestamp: %v", err)
	}

	untilTimestamp := "2023-06-02T15:00:00+00:00"
	parsedUntilTime := time.Time{}
	err = parsedUntilTime.UnmarshalText([]byte(untilTimestamp))
	if err != nil {
		t.Errorf("Error scanning timestamp: %v", err)
	}

	request := httptest.NewRequest("GET", "/045010?from="+url.QueryEscape(fromTimestamp)+"&until="+url.QueryEscape(untilTimestamp), nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)
	router.ServeHTTP(responseRecorder, request)

	var timeseries []types.DataPoint
	err = json.NewDecoder(responseRecorder.Result().Body).Decode(&timeseries)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
		fmt.Fprintf(os.Stdout, "%s\n", responseRecorder.Body)
	}

	if len(timeseries) == 0 {
		t.Errorf("Expected filled timeseries, got %v", timeseries)
	}

	firstPoint := timeseries[0]
	lastPoint := timeseries[len(timeseries)-1]

	if firstPoint.Timestamp.Time.Before(parsedFromTime) {
		t.Errorf("Expected first point timestamp to be before parsedFromTime (%v), got %v", parsedFromTime, firstPoint.Timestamp.Time)
	}

	if lastPoint.Timestamp.Time.After(parsedUntilTime) {
		t.Errorf("Expected last point timestamp to be before parsedUntilTime (%v), got %v", parsedUntilTime, lastPoint.Timestamp.Time)
	}
}
