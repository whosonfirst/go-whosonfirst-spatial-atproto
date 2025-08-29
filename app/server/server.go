package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	gohttp "net/http"

	"github.com/aaronland/go-http/v3/server"
	wof "github.com/whosonfirst/go-whosonfirst-spatial-atproto/http/xrpc/org/whosonfirst"
	app "github.com/whosonfirst/go-whosonfirst-spatial/application"
)

func Run(ctx context.Context) error {

	fs, err := DefaultFlagSet()

	if err != nil {
		return fmt.Errorf("Failed to derive default flag set, %w", err)
	}

	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flag set, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	spatial_opts := &app.SpatialApplicationOptions{
		SpatialDatabaseURI: opts.SpatialDatabaseURI,
		// PropertiesReaderURI:    opts.PropertiesReaderURI,
		// EnableCustomPlacetypes: opts.EnableCustomPlacetypes,
		// CustomPlacetypes:       opts.CustomPlacetypes,
	}

	spatial_app, err := app.NewSpatialApplication(ctx, spatial_opts)

	if err != nil {
		return fmt.Errorf("Failed to create new spatial application, %w", err)
	}

	mux := gohttp.NewServeMux()

	// point-in-polygon handler

	wof_pip_opts := &wof.PointInPolygonHandlerOptions{}

	wof_pip_handler, err := wof.PointInPolygonHandler(spatial_app, wof_pip_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	logger.Debug("Enable point in polygon handler", "endpoint", wof.XRPC_POINT_IN_POLYGON)
	mux.Handle(wof.XRPC_POINT_IN_POLYGON, wof_pip_handler)

	// point-in-polygon from tile handler

	wof_piptile_opts := &wof.PointInPolygonTileHandlerOptions{}

	wof_piptile_handler, err := wof.PointInPolygonTileHandler(spatial_app, wof_piptile_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon from tile handler because %s", err)
	}

	logger.Debug("Enable point in polygon from tile handler", "endpoint", wof.XRPC_POINT_IN_POLYGON_TILE)
	mux.Handle(wof.XRPC_POINT_IN_POLYGON_TILE, wof_piptile_handler)

	// intersects

	wof_intersects_opts := &wof.IntersectsHandlerOptions{}

	wof_intersects_handler, err := wof.IntersectsHandler(spatial_app, wof_intersects_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	logger.Debug("Enable point in polygon handler", "endpoint", wof.XRPC_INTERSECTS)
	mux.Handle(wof.XRPC_INTERSECTS, wof_intersects_handler)

	// record handler

	wof_record_opts := &wof.GetRecordHandlerOptions{}

	wof_record_handler, err := wof.GetRecordHandler(spatial_app, wof_record_opts)

	if err != nil {
		return fmt.Errorf("failed to create get record handler because %s", err)
	}

	logger.Debug("Enable get record handler", "endpoint", wof.XRPC_RECORD)
	mux.Handle(wof.XRPC_RECORD, wof_record_handler)

	// Geocode handler

	wof_geocode_opts := &wof.GeocodeHandlerOptions{
		PlaceholderEndpoint: opts.PlaceholderEndpoint,
	}

	wof_geocode_handler, err := wof.GeocodeHandler(wof_geocode_opts)

	if err != nil {
		return fmt.Errorf("failed to create get geocode handler because %s", err)
	}

	logger.Debug("Enable geocode handler", "endpoint", wof.XRPC_GEOCODE)
	mux.Handle(wof.XRPC_GEOCODE, wof_geocode_handler)

	// Start server

	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server for '%s', %v", server_uri, err)
	}

	logger.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %v", err)
	}

	return nil
}
