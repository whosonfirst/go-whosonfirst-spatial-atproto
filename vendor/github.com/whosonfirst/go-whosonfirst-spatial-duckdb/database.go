package duckdb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/marcboeker/go-duckdb/v2"

	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

type DuckDBSpatialDatabase struct {
	database.SpatialDatabase
	conn         *sql.DB
	database_uri string
}

func init() {
	ctx := context.Background()
	database.RegisterSpatialDatabase(ctx, "duckdb", NewDuckDBSpatialDatabase)
}

func NewDuckDBSpatialDatabase(ctx context.Context, uri string) (database.SpatialDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	database_uri := q.Get("uri")

	conn, err := sql.Open("duckdb", "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create database connection, %w", err)
	}

	extensions := []string{
		"spatial",
		"httpfs",
	}

	for _, ext := range extensions {

		cmds := []string{
			fmt.Sprintf("INSTALL %s", ext),
			fmt.Sprintf("LOAD %s", ext),
		}

		for _, q := range cmds {

			_, err := conn.ExecContext(ctx, q)

			if err != nil {
				return nil, fmt.Errorf("Failed to configure data - query failed, %w (%s)", err, q)
			}
		}
	}

	db := &DuckDBSpatialDatabase{
		conn:         conn,
		database_uri: database_uri,
	}

	return db, nil
}
