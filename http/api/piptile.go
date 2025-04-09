package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aaronland/go-http-sanitize"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/clip"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
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

		// To do: For each result:
		// Fetch geojson Feature
		// Trim/clip geometries to maptile
		// Return GeoJSON

		fc := geojson.NewFeatureCollection()

		for _, r := range intersects_rsp.Results() {

			id, err := strconv.ParseInt(r.Id(), 10, 64)

			if err != nil {
				logger.Error("Failed to derive WOF ID", "id", r.Id(), "error", err)
				xrpcError(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}
			
			body, err := wof_reader.LoadBytes(ctx, app.SpatialDatabase, id)

			if err != nil {
				logger.Error("Failed to load body for WOF ID", "id", r.Id(), "error", err)
				xrpcError(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}
			
			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				logger.Error("Failed to unmarshal feature for WOF ID", "id", r.Id(), "error", err)
				xrpcError(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Clipping happens below
			fc.Append(f)
		}

		col := make([]orb.Geometry, len(fc.Features))

		for idx, f := range fc.Features {
			col[idx] = f.Geometry
		}

		col = clip.Collection(geom.Geometry().Bound(), col)

		for idx, clipped_geom := range col {
			fc.Features[idx].Geometry = clipped_geom
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
