package routes

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/chi/v5"
	wisdomTypes "github.com/wisdom-oss/commonTypes/v2"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"

	"github.com/wisdom-oss/service-smartmeter-rest/globals"
	"github.com/wisdom-oss/service-smartmeter-rest/types"
)

var ErrInvalidTimestamp = wisdomTypes.WISdoMError{
	Type:   "TODO-fill with link to docs",
	Status: 400,
	Title:  "Invalid Timestamp Provided",
	Detail: "A timestamp provided in the request did not follow the required ISO 8691 format",
}

var ErrTimeseriesNotFound = wisdomTypes.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "Timeseries Not Found",
	Detail: "The timeseries with the supplied smartmeter ID does not exist",
}

func TimeSeries(w http.ResponseWriter, r *http.Request) {
	// get the error handlers
	errorHandler := r.Context().Value(middleware.ErrorChannelName).(chan<- interface{})
	statusChannel := r.Context().Value(middleware.StatusChannelName).(<-chan bool)

	// get the time series from the url
	dataSeriesID := chi.URLParam(r, "data-series-id")

	// now check if a timeseries for the smart meter exists
	rawQuery, err := globals.SqlQueries.Raw("timeseries-exists")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var timeseriesExists []bool
	err = pgxscan.Select(r.Context(), globals.Db, &timeseriesExists, rawQuery, dataSeriesID)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	if !timeseriesExists[0] {
		errorHandler <- ErrTimeseriesNotFound
		<-statusChannel
		return
	}

	// since the timeseries exists, now handle parsing the from and until
	// parameter
	// get the query parameter 'from' and 'until'
	fromString, fromSet := r.URL.Query()["from"]
	untilString, untilSet := r.URL.Query()["until"]

	var from, until time.Time

	if fromSet {
		from, err = time.Parse(time.RFC3339, fromString[0])
		if err != nil {
			errorHandler <- ErrInvalidTimestamp
			<-statusChannel
			return
		}
	}

	if untilSet {
		until, err = time.Parse(time.RFC3339, untilString[0])
		if err != nil {
			errorHandler <- ErrInvalidTimestamp
			<-statusChannel
			return
		}
	}

	parameter := []interface{}{dataSeriesID}

	switch {
	// 0 0
	case !fromSet && !untilSet:
		rawQuery, err = globals.SqlQueries.Raw("timeseries")
		break
	// 0 1
	case !fromSet && untilSet:
		rawQuery, err = globals.SqlQueries.Raw("timeseries-daterange-until")
		parameter = append(parameter, until)
		break
	case fromSet && !untilSet:
		rawQuery, err = globals.SqlQueries.Raw("timeseries-daterange-from")
		parameter = append(parameter, from)
		break
	case fromSet && untilSet:
		rawQuery, err = globals.SqlQueries.Raw("timeseries-daterange")
		parameter = append(parameter, from, until)
		break
	}

	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	rows, err := globals.Db.Query(r.Context(), rawQuery, parameter...)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var timeseries []types.DataPoint
	err = pgxscan.ScanAll(&timeseries, rows)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	// since the query has been successful, now determine the output content
	// type
	var outputFormat types.OutputFormat
	outputFormat.FromString(r.URL.Query().Get("format"))

	switch outputFormat {
	case types.CSV:
		w.Header().Set("Content-Type", "text/csv")
		csvWriter := csv.NewWriter(w)
		err = csvWriter.Write([]string{"timestamp", "value"})
		if err != nil {
			errorHandler <- err
			<-statusChannel
			return
		}
		csvWriter.Flush()
		for _, dataPoint := range timeseries {
			timestampString := dataPoint.Timestamp.Time.Format(time.RFC3339)
			value := fmt.Sprintf("%f", dataPoint.Value.Float64)
			err = csvWriter.Write([]string{timestampString, value})
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			csvWriter.Flush()
		}
		break
	case types.JSON:
		// handle JSON output
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(timeseries)
		break

	case types.CBOR:
		w.Header().Set("Content-Type", "application/cbor")
		err = cbor.NewEncoder(w).Encode(timeseries)
		break

	}

}
