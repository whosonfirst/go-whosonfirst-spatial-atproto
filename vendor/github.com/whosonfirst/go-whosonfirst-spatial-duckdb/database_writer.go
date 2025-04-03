package duckdb

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/whosonfirst/go-writer/v3"
)

func init() {
	err := writer.RegisterWriter(context.Background(), "duckd", NewDuckDBSpatialDatabaseWriter)

	if err != nil {
		panic(err)
	}
}

func NewDuckDBSpatialDatabaseWriter(ctx context.Context, uri string) (writer.Writer, error) {
	return NewDuckDBSpatialDatabase(ctx, uri)
}

// Write implements the whosonfirst/go-writer interface so that the database itself can be used as a
// writer.Writer instance (by invoking the `IndexFeature` method).
func (r *DuckDBSpatialDatabase) Write(ctx context.Context, key string, fh io.ReadSeeker) (int64, error) {
	return 0, fmt.Errorf("Not implemented")
}

// WriterURI implements the whosonfirst/go-writer interface so that the database itself can be used as a
// writer.Writer instance
func (r *DuckDBSpatialDatabase) WriterURI(ctx context.Context, str_uri string) string {
	return str_uri
}

// Flush implements the whosonfirst/go-writer interface so that the database itself can be used as a
// writer.Writer instance. This method is a no-op and simply returns `nil`.
func (r *DuckDBSpatialDatabase) Flush(ctx context.Context) error {
	return nil
}

// Close implements the whosonfirst/go-writer interface so that the database itself can be used as a
// writer.Writer instance. This method is a no-op and simply returns `nil`.
func (r *DuckDBSpatialDatabase) Close(ctx context.Context) error {
	return nil
}

// SetLogger implements the whosonfirst/go-writer interface so that the database itself can be used as a
// writer.Writer instance. This method is a no-op and simply returns `nil`.
func (r *DuckDBSpatialDatabase) SetLogger(ctx context.Context, logger *log.Logger) error {
	return nil
}
