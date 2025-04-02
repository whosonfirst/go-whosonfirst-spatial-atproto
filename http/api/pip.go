package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	// "github.com/aaronland/go-http-sanitize"
	// "github.com/whosonfirst/go-whosonfirst-spatial"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

type PointInPolygonHandlerOptions struct {
}

func PointInPolygonHandler(app *spatial_app.SpatialApplication, opts *PointInPolygonHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

		pip_fn, err := query.NewSpatialFunction(ctx, "pip://")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		// START OF read from params...
		
		var pip_query *query.SpatialQuery

		dec := json.NewDecoder(req.Body)
		err = dec.Decode(&pip_query)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		// END OF read from params...
		
		pip_rsp, err := query.ExecuteQuery(ctx, app.SpatialDatabase, pip_fn, pip_query)

		if err != nil {
			logger.Error("Failed to execute point in polygon query", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		/*
		if len(pip_query.Properties) > 0 {

			props_opts := &spatial.PropertiesResponseOptions{
				Reader:       app.PropertiesReader,
				Keys:         pip_query.Properties,
				SourcePrefix: "properties",
			}

			app.Monitor.Signal(ctx, timings.SinceStart, timingsPIPProperties)

			props_rsp, err := spatial.PropertiesResponseResultsWithStandardPlacesResults(ctx, props_opts, pip_rsp)

			app.Monitor.Signal(ctx, timings.SinceStop, timingsPIPProperties)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			enc := json.NewEncoder(rsp)
			err = enc.Encode(props_rsp)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}
		*/

		// Generate ATGeo here...
		
		enc := json.NewEncoder(rsp)
		err = enc.Encode(pip_rsp)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	pip_handler := http.HandlerFunc(fn)
	return pip_handler, nil
}
