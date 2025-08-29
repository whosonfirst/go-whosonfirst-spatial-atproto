package duckdb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-reader/v2"
	"github.com/whosonfirst/go-whosonfirst-format"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

func init() {

	err := reader.RegisterReader(context.Background(), "duckd", NewDuckDBSpatialDatabaseReader)

	if err != nil {
		panic(err)
	}
}

func NewDuckDBSpatialDatabaseReader(ctx context.Context, uri string) (reader.Reader, error) {
	return NewDuckDBSpatialDatabase(ctx, uri)
}

// Read implements the whosonfirst/go-reader interface so that the database itself can be used as a
// reader.Reader instance (reading features from the `geojson` table.
func (db *DuckDBSpatialDatabase) Read(ctx context.Context, str_uri string) (io.ReadSeekCloser, error) {

	// return nil, fmt.Errorf("Not implemented")

	id, _, err := uri.ParseURI(str_uri)

	if err != nil {
		return nil, err
	}

	// There is a problem here that I am just not seeing right now...

	q := fmt.Sprintf(`SELECT parent_id, name, placetype, country, repo, lat, lon, modified,  ST_AsText(ST_GeomFromWKB(geometry)) AS wkt_geom FROM read_parquet('%s') WHERE id = ?`, db.database_uri)

	row := db.conn.QueryRowContext(ctx, q, id)

	var parent_id int64
	var name string
	var placetype string
	var country string
	var repo string
	var lat float64
	var lon float64
	var str_lastmod string
	var wkt_geom string

	err = row.Scan(&parent_id, &name, &placetype, &country, &repo, &lat, &lon, &str_lastmod, &wkt_geom)

	if err != nil {
		return nil, fmt.Errorf("Failed to scan row for %d, %w", id, err)
	}

	lastmod := int64(0)

	t, err := time.Parse(time.RFC3339, str_lastmod)

	if err != nil {
		slog.Warn("Failed to parse lastmod string", "lastmod", str_lastmod, "error", err)
	} else {
		lastmod = t.Unix()
	}

	orb_geom, err := wkt.Unmarshal(wkt_geom)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal geometry, %w", err)
	}

	f := geojson.NewFeature(orb_geom)

	f.Properties = map[string]interface{}{
		"wof:id":           id,
		"wof:name": name,
		"wof:parent_id":    parent_id,
		"wof:placetype":    placetype,
		"wof:country":      country,
		"wof:repo":         repo,
		"geom:latitude":    lat,
		"geom:longitude":   lon,
		"wof:lastmodified": lastmod,
	}

	f.ID = id

	enc_f, err := format.FormatFeature(f)

	if err != nil {
		return nil, fmt.Errorf("Failed to format feature, %w", err)
	}

	br := bytes.NewReader(enc_f)
	return ioutil.NewReadSeekCloser(br)
}

// Exists returns a boolean value indicating whether 'str_uri' already exists in the database.
func (db *DuckDBSpatialDatabase) Exists(ctx context.Context, str_uri string) (bool, error) {


	id, _, err := uri.ParseURI(str_uri)

	if err != nil {
		return false, err
	}

	q := fmt.Sprintf("SELECT 1 FROM read_parquet('%s') WHERE id = ?", db.database_uri)

	row := db.conn.QueryRowContext(ctx, q, id)

	var one int

	err = row.Scan(&one)

	if err != nil {

		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil	
}

// ReadURI implements the whosonfirst/go-reader interface so that the database itself can be used as a
// reader.Reader instance
func (db *DuckDBSpatialDatabase) ReaderURI(ctx context.Context, str_uri string) string {
	return str_uri
}
