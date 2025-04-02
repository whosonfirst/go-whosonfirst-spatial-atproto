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

type IntersectsHandlerOptions struct {
}

func IntersectsHandler(app *spatial_app.SpatialApplication, opts *IntersectsHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

		if req.Method != "POST" {
			http.Error(rsp, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		intersects_fn, err := query.NewSpatialFunction(ctx, "intersects://")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		// START OF ...
		
		var intersects_query *query.SpatialQuery

		dec := json.NewDecoder(req.Body)
		err = dec.Decode(&intersects_query)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		// END OF ...
		
		intersects_rsp, err := query.ExecuteQuery(ctx, app.SpatialDatabase, intersects_fn, intersects_query)

		if err != nil {
			logger.Error("Failed to execute intersects query", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		/*
		if len(intersects_query.Properties) > 0 {

			props_opts := &spatial.PropertiesResponseOptions{
				Reader:       app.PropertiesReader,
				Keys:         intersects_query.Properties,
				SourcePrefix: "properties",
			}

			app.Monitor.Signal(ctx, timings.SinceStart, timingsIntersectsProperties)

			props_rsp, err := spatial.PropertiesResponseResultsWithStandardPlacesResults(ctx, props_opts, intersects_rsp)

			app.Monitor.Signal(ctx, timings.SinceStop, timingsIntersectsProperties)

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
		
		enc := json.NewEncoder(rsp)
		err = enc.Encode(intersects_rsp)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	intersects_handler := http.HandlerFunc(fn)
	return intersects_handler, nil
}
