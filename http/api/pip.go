package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http/v3/sanitize"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

type PointInPolygonHandlerOptions struct{}

func PointInPolygonHandler(app *spatial_app.SpatialApplication, opts *PointInPolygonHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

		pip_fn, err := query.NewSpatialFunction(ctx, "pip://")

		if err != nil {
			logger.Error("Failed to construct spatial fuction (pip://)", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		lat, err := sanitize.GetFloat64(req, P_LATITUDE)

		if err != nil {
			logger.Error("Failed to derive latitude", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		lon, err := sanitize.GetFloat64(req, P_LONGITUDE)

		if err != nil {
			logger.Error("Failed to derive longitude", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF put me in a function...

		pt := orb.Point([2]float64{lon, lat})
		geom := geojson.NewGeometry(pt)

		pip_query := &query.SpatialQuery{
			Geometry: geom,
		}

		// END OF put me in a function...

		// Something something something this information is not (?) in the parquet files...
		// pip_query.IsCurrent = []int64{1}

		pip_rsp, err := query.ExecuteQuery(ctx, app.SpatialDatabase, pip_fn, pip_query)

		if err != nil {
			logger.Error("Failed to execute point in polygon query", "error", err)
			xrpcError(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Generate ATGeo here...

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(pip_rsp)

		if err != nil {
			xrpcError(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	pip_handler := http.HandlerFunc(fn)
	return pip_handler, nil
}
