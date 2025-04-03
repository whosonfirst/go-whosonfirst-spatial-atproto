# go-whosonfirst-spatial-duckdb

Go package to implement the [whosonfirst/go-whosonfirst-spatial](https://github.com/whosonfirst/go-whosonfirst-spatial) interfaces using a GeoParquet database (and DuckDB).

## Documentation

Documentation is incomplete at this time.

## Work in progress

This package should still be considered "work in progress". The basics all seem to work but, in addition to any documentation which still needs to be written, some bits of functionality may still be absent.

## Building

Because of their size the various DuckDB `libduckdb.a` files are explicitly excluded from the `vendor` directory in this package. Per the [vendoring documentation](https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring) in the `marcboeker/go-duckdb` package the easiest thing to do is this:

```
$> go install github.com/goware/modvendor@latest
$> go mod vendor
$> modvendor -copy="**/*.a **/*.h" -v
```

This package includes handy `modvendor` Makefile target for automating some of that. For example:

```
$> make modvendor
modvendor -copy="**/*.a **/*.h" -v
vendoring github.com/apache/arrow-go/v18/internal/utils/_lib/arch.h
vendoring github.com/apache/arrow-go/v18/arrow/cdata/utils.h
vendoring github.com/apache/arrow-go/v18/arrow/cdata/arrow/c/helpers.h
vendoring github.com/apache/arrow-go/v18/arrow/memory/internal/cgoalloc/allocator.h
vendoring github.com/apache/arrow-go/v18/arrow/compute/internal/kernels/_lib/types.h
vendoring github.com/apache/arrow-go/v18/arrow/compute/internal/kernels/_lib/vendored/safe-math.h
vendoring github.com/apache/arrow-go/v18/arrow/math/_lib/arch.h
vendoring github.com/apache/arrow-go/v18/arrow/memory/_lib/arch.h
vendoring github.com/apache/arrow-go/v18/arrow/cdata/arrow/c/abi.h
vendoring github.com/apache/arrow-go/v18/arrow/memory/internal/cgoalloc/helpers.h
vendoring github.com/google/flatbuffers/goldens/cpp/basic_generated.h
vendoring github.com/marcboeker/go-duckdb/deps/windows_amd64/libduckdb.a
vendoring github.com/marcboeker/go-duckdb/deps/linux_amd64/libduckdb.a
vendoring github.com/marcboeker/go-duckdb/deps/linux_arm64/libduckdb.a
vendoring github.com/marcboeker/go-duckdb/duckdb.h
vendoring github.com/marcboeker/go-duckdb/deps/darwin_arm64/libduckdb.a
vendoring github.com/marcboeker/go-duckdb/deps/freebsd_amd64/libduckdb.a
vendoring github.com/marcboeker/go-duckdb/deps/darwin_amd64/libduckdb.a
vendoring golang.org/x/tools/internal/gcimporter/testdata/versions/test_go1.20_u.a
```

It's an unfortunate extra set of steps but "oh well".

_This package has not been tested with `marcboeker/go-duckdb/v2` yet._

## GeoParquet files

As of this writing this package targets the GeoParquet files [produced by geocode.earth](https://geocode.earth/data/whosonfirst/combined/).

## Spatial database URIs

Spatial database URIs for DuckDB/GeoParquet databases take the form of:

```
duckdb://?{QUERY_PARAMETERS}
```

### Query parameters

| Name | Value | Required | Notes |
| --- | --- | --- | --- |
| uri | string | yes | A valid string that can be passed to the DuckDB `read_parquet` command. |

For example:

```
duckdb://?uri=/usr/local/whosonfirst/go-whosonfirst-spatial-duckdb/fixtures/sf_county.parquet
```

## Example

```
import (
       "context"
       "fmt"
       
       "github.com/paulmach/orb"       
       "github.com/whosonfirst/go-whosonfirst-spatial/database"
       _ "github.com/whosonfirst/go-whosonfirst-spatial-duckdb"       
)

func main(){

     ctx := context.Background()

     db, _ := database.NewSpatialDatabase(ctx, "duckdb://?uri=/usr/local/whosonfirst/go-whosonfirst-spatial-duckdb/fixtures/sf_county.parquet")

     lat := 37.621131
     lon := -122.384292
     
     pt := orb.Point{ lon, lat }
     
     spr, _ := db.PointInPolygon(ctx, pt)

     for _, r := range spr.Results(){
     	fmt.Printf("%s %s\n", r.Id(), r.Name())
     }
}
```

## Tools

```
$> make cli
go build -mod readonly -ldflags="-s -w" -o bin/http-server cmd/http-server/main.go
go build -mod readonly -ldflags="-s -w" -o bin/pip cmd/pip/main.go
go build -mod readonly -ldflags="-s -w" -o bin/intersects cmd/intersects/main.go
```

### pip

Documentation for the `pip` tool can be found in [cmd/pip/README.md](cmd/pip/README.md).

### intersects

Documentation for the `intersects` tool can be found in [cmd/intersects/README.md](cmd/intersects/README.md).

### http-server

Documentation for the `http-server` tool can be found in [cmd/http-server/README.md](cmd/http-server/README.md).

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/marcboeker/go-duckdb
* https://duckdb.org/
* https://geocode.earth/data/whosonfirst/combined/