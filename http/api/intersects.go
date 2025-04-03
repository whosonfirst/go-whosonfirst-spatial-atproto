package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geojson"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

type IntersectsHandlerOptions struct{}

func IntersectsHandler(app *spatial_app.SpatialApplication, opts *IntersectsHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

		intersects_fn, err := query.NewSpatialFunction(ctx, "intersects://")

		if err != nil {
			logger.Error("Failed to construct spatial fuction (intersects://)", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF read from params...

		str_wkt, err := sanitize.GetString(req, P_GEOMETRY)

		if err != nil {
			logger.Error("Failed to derive geometry", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF put me in a function...

		orb_geom, err := wkt.Unmarshal(str_wkt)

		if err != nil {
			logger.Error("Failed to unmarshal geometry", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		geom := geojson.NewGeometry(orb_geom)

		intersects_query := &query.SpatialQuery{
			Geometry: geom,
		}

		// END OF put me in a function...

		// intersects_query.IsCurrent = []int64{1}

		intersects_rsp, err := query.ExecuteQuery(ctx, app.SpatialDatabase, intersects_fn, intersects_query)

		if err != nil {
			logger.Error("Failed to execute intersects query", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

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
