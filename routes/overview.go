package routes

import (
	"encoding/json"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	middleware "github.com/wisdom-oss/microservice-middlewares/v4"

	"github.com/wisdom-oss/service-smartmeter-rest/globals"
	"github.com/wisdom-oss/service-smartmeter-rest/types"
)

func Overview(w http.ResponseWriter, r *http.Request) {
	// get the error handlers
	errorHandler := r.Context().Value(middleware.ErrorChannelName).(chan<- interface{})
	statusChannel := r.Context().Value(middleware.StatusChannelName).(<-chan bool)

	// now get the raw sql query from the file
	rawQuery, err := globals.SqlQueries.Raw("timeseries-overview")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	rows, err := globals.Db.Query(r.Context(), rawQuery)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var overviews []types.TimeseriesInformation
	err = pgxscan.ScanAll(&overviews, rows)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(overviews)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}
}
