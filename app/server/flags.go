package server

import (
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
	// "github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

var verbose bool
var server_uri string
var spatial_database_uri string
var placeholder_endpoint string

// var properties_reader_uri string
// var enable_custom_placetypes bool
// var custom_placetypes string

func DefaultFlagSet() (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("server")

	available_databases := database.Schemes()
	desc_databases := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-spatial/data.SpatialDatabase URI. options are: %s", available_databases)

	fs.StringVar(&spatial_database_uri, "spatial-database-uri", "rtree://", desc_databases)

	fs.StringVar(&placeholder_endpoint, "placeholder-endpoint", "http://localhost:3000", "The URL that the Placeholder API is listening on.")

	// available_readers := reader.Schemes()
	// desc_readers := fmt.Sprintf("A valid whosonfirst/go-reader.Reader URI. Available options are: %s", available_readers)

	// fs.StringVar(&properties_reader_uri, "properties-reader-uri", "", fmt.Sprintf("%s. If the value is {spatial-database-uri} then the value of the '-spatial-database-uri' implements the reader.Reader interface and will be used.", desc_readers))

	// fs.BoolVar(&enable_custom_placetypes, "enable-custom-placetypes", false, "Enable wof:placetype values that are not explicitly defined in the whosonfirst/go-whosonfirst-placetypes repository.")

	// fs.StringVar(&custom_placetypes, "custom-placetypes", "", "A JSON-encoded string containing custom placetypes defined using the syntax described in the whosonfirst/go-whosonfirst-placetypes repository.")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs, nil
}
