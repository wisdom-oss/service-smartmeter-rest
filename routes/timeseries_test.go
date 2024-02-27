package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
	validator "openapi.tanna.dev/go/validator/openapi3"
)

func TestTimeSeries_InvalidDataSeriesID(t *testing.T) {
	contract, err := openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		t.Errorf("unable to load openapi contract: %s", err.Error())
	}
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/0", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d\n%s", http.StatusNotFound, responseRecorder.Code, responseRecorder.Body.String())
	}
}

func TestTimeSeries_JSON(t *testing.T) {
	contract, err := openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		t.Errorf("unable to load openapi contract: %s", err.Error())
	}
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}

func TestTimeSeries_CSV(t *testing.T) {
	contract, err := openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		t.Errorf("unable to load openapi contract: %s", err.Error())
	}
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/{data-series-id}", TimeSeries)

	request := httptest.NewRequest("GET", "/045010?format=csv", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}
