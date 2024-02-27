package routes

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/qustavo/dotsql"
	wisdomType "github.com/wisdom-oss/commonTypes"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
	validator "openapi.tanna.dev/go/validator/openapi3"

	"github.com/wisdom-oss/service-smartmeter-rest/globals"
)

func TestMain(m *testing.M) {
	// load the variables found in the .env file into the process environment
	godotenv.Load()
	location, locationChanged := os.LookupEnv("ENV_CONFIG_LOCATION")
	if !locationChanged {
		location = "./environment.json"
	}
	var c wisdomType.EnvironmentConfiguration
	c.PopulateFromFilePath(location)
	globals.Environment, _ = c.ParseEnvironment()
	address := fmt.Sprintf("postgres://%s:%s@%s:%s/wisdom",
		globals.Environment["PG_USER"], globals.Environment["PG_PASS"],
		globals.Environment["PG_HOST"], globals.Environment["PG_PORT"])

	config, _ := pgxpool.ParseConfig(address)
	globals.Db, _ = pgxpool.NewWithConfig(context.Background(), config)
	globals.SqlQueries, _ = dotsql.LoadFromFile(globals.Environment["QUERY_FILE_LOCATION"])
	os.Exit(m.Run())
}

func TestOverview_JSON(t *testing.T) {
	contract, err := openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		t.Errorf("unable to load openapi contract: %s", err.Error())
	}
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", Overview)

	request := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}
