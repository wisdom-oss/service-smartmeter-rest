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
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"
	validator "openapi.tanna.dev/go/validator/openapi3"

	"github.com/wisdom-oss/service-smartmeter-rest/globals"
)

var contract *openapi3.T

const responseFormatString = "Response Code: %d\nResponse Status:%s\nContent-Length: %s\nContent-Type: %s\n\n%s"

func TestMain(m *testing.M) {
	// load the variables found in the .env file into the process environment
	godotenv.Load()
	var c wisdomType.EnvironmentConfiguration
	c.PopulateFromFilePath("./resources/environment.json")
	globals.Environment, _ = c.ParseEnvironment()
	address := fmt.Sprintf("postgres://%s:%s@%s:%s/wisdom",
		globals.Environment["PG_USER"], globals.Environment["PG_PASS"],
		globals.Environment["PG_HOST"], globals.Environment["PG_PORT"])

	config, _ := pgxpool.ParseConfig(address)
	globals.Db, _ = pgxpool.NewWithConfig(context.Background(), config)
	globals.SqlQueries, _ = dotsql.LoadFromFile("./resources/queries.sql")

	var err error
	contract, err = openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestOverview_JSON(t *testing.T) {
	router := chi.NewRouter()
	router.Use(middleware.ErrorHandler)
	router.Get("/", Overview)

	request := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	_ = validator.NewValidator(contract).ForTest(t, responseRecorder, request)

	router.ServeHTTP(responseRecorder, request)
}
