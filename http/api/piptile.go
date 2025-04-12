package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	orb_maptile "github.com/paulmach/orb/maptile"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
	"github.com/whosonfirst/go-whosonfirst-spatial/maptile"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

type PointInPolygonTileHandlerOptions struct{}

func PointInPolygonTileHandler(app *spatial_app.SpatialApplication, opts *PointInPolygonTileHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

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

		zm := orb_maptile.Zoom(uint32(z))
		map_t := orb_maptile.New(uint32(x), uint32(y), zm)

		// Something something something add filters here...
		spatial_q := &query.SpatialQuery{}

		fc, err := maptile.PointInPolygonCandidateFeaturesFromTile(ctx, app.SpatialDatabase, spatial_q, &map_t)

		if err != nil {
			logger.Error("Failed to execute point in polygon query", "error", err)
			xrpcError(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Generate ATGeo here...

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(fc)

		if err != nil {
			xrpcError(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	piptile_handler := http.HandlerFunc(fn)
	return piptile_handler, nil
}
