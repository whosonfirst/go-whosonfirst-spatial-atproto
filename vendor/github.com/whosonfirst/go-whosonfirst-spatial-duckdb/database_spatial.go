package duckdb

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"log/slog"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

func (db *DuckDBSpatialDatabase) IndexFeature(context.Context, []byte) error {
	return fmt.Errorf("Not implemented.")
}

func (db *DuckDBSpatialDatabase) RemoveFeature(context.Context, string) error {
	return fmt.Errorf("Not implemented.")
}

func (db *DuckDBSpatialDatabase) PointInPolygon(ctx context.Context, coord *orb.Point, filters ...spatial.Filter) (spr.StandardPlacesResults, error) {

	results := make([]spr.StandardPlacesResult, 0)
	var err error

	for wof_spr, pip_err := range db.PointInPolygonWithIterator(ctx, coord, filters...) {

		if err != nil {
			err = pip_err
			break
		}

		if wof_spr != nil {
			results = append(results, wof_spr)
		}
	}

	if err != nil {
		return nil, err
	}

	return NewSPRResults(results), nil
}

func (db *DuckDBSpatialDatabase) Intersects(ctx context.Context, poly orb.Geometry, filters ...spatial.Filter) (spr.StandardPlacesResults, error) {

	results := make([]spr.StandardPlacesResult, 0)
	var err error

	for wof_spr, pip_err := range db.IntersectsWithIterator(ctx, poly, filters...) {

		if err != nil {
			err = pip_err
			break
		}

		results = append(results, wof_spr)
	}

	if err != nil {
		return nil, err
	}

	return NewSPRResults(results), nil
}

func (db *DuckDBSpatialDatabase) PointInPolygonWithIterator(ctx context.Context, coord *orb.Point, filters ...spatial.Filter) iter.Seq2[spr.StandardPlacesResult, error] {

	logger := slog.Default()
	logger = logger.With("coordinate", coord)

	return func(yield func(spr.StandardPlacesResult, error) bool) {

		yield(nil, fmt.Errorf("Not implemented."))

		// See what's happening here? We are violating the cardinal rule of SQL
		// wrangling by NOT passing in variables to be quoted or escaped as necessary
		// by the SQL driver. That's because every time I do I get errors like this:
		// 2024/12/12 10:46:44 ERROR Q error="Binder Error: Referenced column \"POINT(? ?)\" not found in FROM clause!\nCandidate bindings: \"read_parquet.id\"\nLINE 1: ... WHERE ST_Contains(geometry::GEOMETRY, \"POINT(? ?)\"::GEOMETRY)\n

		q := fmt.Sprintf(`SELECT id, parent_id, name, placetype, country, repo, lat, lon, min_lat, min_lon, max_lat, max_lon, modified FROM read_parquet('%s') WHERE ST_Contains(ST_GeomFromWKB(geometry), 'POINT(%f %f)')`, db.database_uri, coord.X(), coord.Y())

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			logger.Error("Query failed", "error", err)
			yield(nil, err)
			return
		}

		defer rows.Close()

		for r, err := range db.rowsToSPR(ctx, rows, filters...) {

			yield(r, err)

			if err != nil {
				return
			}
		}
	}
}

func (db *DuckDBSpatialDatabase) IntersectsWithIterator(ctx context.Context, geom orb.Geometry, filters ...spatial.Filter) iter.Seq2[spr.StandardPlacesResult, error] {

	geom_type := geom.GeoJSONType()

	logger := slog.Default()
	logger = logger.With("geometry", geom_type)

	return func(yield func(spr.StandardPlacesResult, error) bool) {

		var wkt_str string

		switch geom_type {
		case "Polygon", "MulitPolygon":
			wkt_str = wkt.MarshalString(geom)
		default:
			logger.Error("Invalid geometry type")
			yield(nil, fmt.Errorf("Invalid geometry type"))
			return
		}

		q := fmt.Sprintf(`SELECT id, parent_id, name, placetype, country, repo, lat, lon, min_lat, min_lon, max_lat, max_lon, modified FROM read_parquet('%s') WHERE ST_Intersects(ST_GeomFromWKB(geometry), '%s')`, db.database_uri, wkt_str)

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			logger.Error("Query failed", "error", err)
			yield(nil, err)
			return
		}

		defer rows.Close()
		
		for r, err := range db.rowsToSPR(ctx, rows, filters...) {

			if err != nil {
				logger.Error("Failed to derive SPR from row", "error", err)
				return
			}
			
			if !yield(r, err){
				break
			}

		}
	}
}

func (db *DuckDBSpatialDatabase) Disconnect(ctx context.Context) error {
	return db.conn.Close()
}

func (db *DuckDBSpatialDatabase) rowsToSPR(ctx context.Context, rows *sql.Rows, filters ...spatial.Filter) iter.Seq2[spr.StandardPlacesResult, error] {

	logger := slog.Default()
	
	return func(yield func(spr.StandardPlacesResult, error) bool) {

		for rows.Next() {

			var id int64
			var parent_id int64
			var name string
			var placetype string
			var country string
			var repo string
			var lat float64
			var lon float64
			var min_lat float64
			var min_lon float64
			var max_lat float64
			var max_lon float64
			var str_lastmod string

			err := rows.Scan(&id, &parent_id, &name, &placetype, &country, &repo, &lat, &lon, &min_lat, &min_lon, &max_lat, &max_lon, &str_lastmod)

			if err != nil {
				logger.Error("Row scanning failed", "error", err)
				yield(nil, err)
				break
			}

			rel_path, err := uri.Id2RelPath(id)

			if err != nil {
				logger.Error("Failed to derive rel path for ID", "id", id, "error", err)
				yield(nil, err)
				break
			}

			lastmod := int64(0)

			t, err := time.Parse(time.RFC3339, str_lastmod)

			if err != nil {
				logger.Warn("Failed to parse lastmod string", "lastmod", str_lastmod, "error", err)
			} else {
				lastmod = t.Unix()
			}

			wof_spr := &spr.WOFStandardPlacesResult{
				WOFId:           id,
				WOFParentId:     parent_id,
				WOFName:         name,
				WOFPlacetype:    placetype,
				WOFCountry:      country,
				WOFRepo:         repo,
				WOFPath:         rel_path,
				MZURI:           rel_path,
				MZLatitude:      lat,
				MZLongitude:     lon,
				MZMinLatitude:   min_lat,
				MZMinLongitude:  min_lon,
				MZMaxLatitude:   max_lat,
				MZMaxLongitude:  max_lon,
				WOFLastModified: lastmod,
				// Things we don't know
				WOFSupersededBy: make([]int64, 0),
				WOFSupersedes:   make([]int64, 0),
				WOFBelongsTo:    make([]int64, 0),
				MZIsCurrent:     int64(-1),
				MZIsCeased:      int64(-1),
				MZIsDeprecated:  int64(-1),
				MZIsSuperseded:  int64(-1),
				MZIsSuperseding: int64(-1),
			}

			filters_ok := true

			for _, f := range filters {

				err = filter.FilterSPR(f, wof_spr)

				if err != nil {
					logger.Debug("Filter error", "error", err)
					filters_ok = false
					break
				}
			}

			if !filters_ok {
				continue
			}

			yield(wof_spr, nil)
		}

		err := rows.Close()

		if err != nil {
			logger.Error("Failed to close rows", "error", err)
			yield(nil, err)
		}
	}
}
