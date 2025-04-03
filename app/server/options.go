package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI              string
	SpatialDatabaseURI     string
	// PropertiesReaderURI    string
	// EnableCustomPlacetypes bool
	// CustomPlacetypes       string
	Verbose                bool
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "WHOSONFIRST")

	if err != nil {
		return nil, fmt.Errorf("Failed to set flags from environment variables, %v", err)
	}

	opts := &RunOptions{
		ServerURI:              server_uri,
		SpatialDatabaseURI:     spatial_database_uri,
		// PropertiesReaderURI:    properties_reader_uri,
		// EnableCustomPlacetypes: enable_custom_placetypes,
		// CustomPlacetypes:       custom_placetypes,
		Verbose:                verbose,
	}

	return opts, nil
}
