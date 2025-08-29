package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	gohttp "net/http"

	"github.com/aaronland/go-http/v3/server"
	"github.com/whosonfirst/go-whosonfirst-spatial-atproto/http/api"
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

	api_pip_opts := &api.PointInPolygonHandlerOptions{}

	api_pip_handler, err := api.PointInPolygonHandler(spatial_app, api_pip_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	logger.Debug("Enable point in polygon handler", "endpoint", api.XRPC_POINT_IN_POLYGON)
	mux.Handle(api.XRPC_POINT_IN_POLYGON, api_pip_handler)

	// point-in-polygon from tile handler

	api_piptile_opts := &api.PointInPolygonTileHandlerOptions{}

	api_piptile_handler, err := api.PointInPolygonTileHandler(spatial_app, api_piptile_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon from tile handler because %s", err)
	}

	logger.Debug("Enable point in polygon from tile handler", "endpoint", api.XRPC_POINT_IN_POLYGON_TILE)
	mux.Handle(api.XRPC_POINT_IN_POLYGON_TILE, api_piptile_handler)

	// intersects

	api_intersects_opts := &api.IntersectsHandlerOptions{}

	api_intersects_handler, err := api.IntersectsHandler(spatial_app, api_intersects_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	logger.Debug("Enable point in polygon handler", "endpoint", api.XRPC_INTERSECTS)
	mux.Handle(api.XRPC_INTERSECTS, api_intersects_handler)

	// record handler

	api_record_opts := &api.GetRecordHandlerOptions{}

	api_record_handler, err := api.GetRecordHandler(spatial_app, api_record_opts)

	if err != nil {
		return fmt.Errorf("failed to create get record handler because %s", err)
	}

	logger.Debug("Enable get record handler", "endpoint", api.XRPC_RECORD)
	mux.Handle(api.XRPC_RECORD, api_record_handler)

	// Geocode handler

	api_geocode_opts := &api.GeocodeHandlerOptions{
		PlaceholderEndpoint: opts.PlaceholderEndpoint,
	}

	api_geocode_handler, err := api.GeocodeHandler(api_geocode_opts)

	if err != nil {
		return fmt.Errorf("failed to create get geocode handler because %s", err)
	}

	logger.Debug("Enable geocode handler", "endpoint", api.XRPC_GEOCODE)
	mux.Handle(api.XRPC_GEOCODE, api_geocode_handler)

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
