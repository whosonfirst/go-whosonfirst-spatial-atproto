package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

type PointInPolygonTileHandlerOptions struct{}

func PointInPolygonTileHandler(app *spatial_app.SpatialApplication, opts *PointInPolygonTileHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

		intersects_fn, err := query.NewSpatialFunction(ctx, "intersects://")

		if err != nil {
			logger.Error("Failed to construct spatial fuction (intersects://)", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		z, err := sanitize.GetInt(req, P_TILE_Z)

		if err != nil {
			logger.Error("Failed to derive z", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		x, err := sanitize.GetInt(req, P_TILE_X)

		if err != nil {
			logger.Error("Failed to derive x", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		y, err := sanitize.GetInt(req, P_TILE_Y)

		if err != nil {
			logger.Error("Failed to derive y", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF put me in a function...

		zm := maptile.Zoom(uint32(z))
		tl := maptile.New(uint32(x), uint32(y), zm)

		bounds := tl.Bound()
		orb_geom := bounds.ToPolygon()

		geom := geojson.NewGeometry(orb_geom)

		intersects_query := &query.SpatialQuery{
			Geometry: geom,
		}

		// END OF put me in a function...

		// Something something something this information is not (?) in the parquet files...
		// intersects_query.IsCurrent = []int64{1}

		intersects_rsp, err := query.ExecuteQuery(ctx, app.SpatialDatabase, intersects_fn, intersects_query)

		if err != nil {
			logger.Error("Failed to execute point in polygon query", "error", err)
			xrpcError(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Generate ATGeo here...

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(intersects_rsp)

		if err != nil {
			xrpcError(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	piptile_handler := http.HandlerFunc(fn)
	return piptile_handler, nil
}
