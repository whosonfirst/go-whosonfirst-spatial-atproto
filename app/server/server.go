package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	gohttp "net/http"

	"github.com/aaronland/go-http-server"
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

	logger := slog.Default()

	spatial_opts := &app.SpatialApplicationOptions{
		SpatialDatabaseURI:     opts.SpatialDatabaseURI,
		PropertiesReaderURI:    opts.PropertiesReaderURI,
		EnableCustomPlacetypes: opts.EnableCustomPlacetypes,
		CustomPlacetypes:       opts.CustomPlacetypes,
	}

	spatial_app, err := app.NewSpatialApplication(ctx, spatial_opts)

	if err != nil {
		return fmt.Errorf("Failed to create new spatial application, %w", err)
	}

	mux := gohttp.NewServeMux()


	// point-in-polygon handlers

	api_pip_opts := &api.PointInPolygonHandlerOptions{}

	api_pip_handler, err := api.PointInPolygonHandler(spatial_app, api_pip_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	mux.Handle("/xrpc/org.whosonfirst.PointInPolygon", api_pip_handler)

	// intersects

	api_intersects_opts := &api.IntersectsHandlerOptions{}

	api_intersects_handler, err := api.IntersectsHandler(spatial_app, api_intersects_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	mux.Handle("/xrpc/org.whosonfirst.Intersects", api_intersects_handler)

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
