# go-whosonfirst-spatial-atproto

## This work is experimental

Nothing about this package should be considered "stable".

It is working code to try and identify the outlines for a series of open questions about how geo services should be exposed in an ATProto context.

## Building

At minimum you will need to install the [Go](https://go.dev/dl) programming language and, ideally, the `make` command line tool to run these tools. If you don't have `make` installed you can copy-paste all the examples, below, in to a terminal.

This code is derived from the [whosonfirst/go-whosonfirst-spatial](https://github.com/whosonfirst/?q=go-whosonfirst-spatial&type=all&language=&sort=) packages which are organized in such a way that interface and service definitions are separate from any given database implementation.

For the sake of brevity this package is a "smushing up" of the [whosonfirst/go-whosonfirst-spatial-www](https://github.com/whosonfirst/go-whosonfirst-spatial-www) and [whosonfirst/go-whosonfirst-spatial-duckdb](https://github.com/whosonfirst/go-whosonfirst-spatial-duckdb) packages. Assuming this code goes anywhere beyond the experimental phase it will likely be broken up in discrete pieces. For now it's just one big messy party.

### Data sources

As mentioned this package defaults to using [DuckDB](https://duckdb.org/) as its database engine. Other databases are supported but they are not enabled by default. Poke me if you want to know how to use a different database (SQLite, PMTiles, etc.)

There is a default DuckDB GeoParquet file for [San Francisco county](https://spelunker.whosonfirst.org/id/102087579) and all its descendants in the [fixtures](fixtures) directory. This file is bundled using Git LFS so depending on your setup you may need to do additional `git clone whatever` commands.

This file was derived from the [Geocode Earth GeoParquet download](https://geocode.earth/data/whosonfirst/combined/) so that will work too if you feel like downloading a 6GB parquet file.

The Geocode Earth GeoParquet files do not contain relevant Who's On First properties for filtering by date or "current-ness" so occasionally places like the [NFL Experience](https://spelunker.whosonfirst.org/id/420564211) microhood, which only lasted a week in 2016, may appear (below).

### DuckDB

Under the hood this file uses the [marcboeker/go-duckdb](https://github.com/marcboeker/go-duckdb) package which bundles all of the various `libduck.*` files. Because of their size those files are explicitly excluded from the `vendor` directory in this package. Per the [vendoring documentation](https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring) in the `marcboeker/go-duckdb` package the easiest thing to do is this:

```
$> go install github.com/goware/modvendor@latest
$> go mod vendor
$> modvendor -copy="**/*.a **/*.h" -v
```

This package includes handy `modvendor` Makefile target for automating some of that.

## Things this package does

* Point-in-polygon spatial queries.
* Intersects (with geometry) spatial queries.
* Point-in-polygon from ZXY map tile spatial queries (incomplete).
* Basic "get record" lookups.

## Things this package doesn't do yet

* Nearby (radial) queries
* Venues. See [whosonfirst/go-whosonfirst-external](https://github.com/whosonfirst/go-whosonfirst-external) and [whosonfirst/whosonfirst-external-duckdb](https://github.com/whosonfirst/whosonfirst-external-duckdb) for possible approaches.
* Geocoding. See [pelias/placeholder](https://github.com/pelias/placeholder/) and [sfomuseum/go-placeholder-client](https://github.com/sfomuseum/go-placeholder-client) for possible approaches.

It probably doesn't return things in the correct ATProto formats. Any pointers suitable for someone-who-is-already-juggling-too-many-things are welcome.

## Example

The easiest way to get started is to run the handy `debug` Makfile target:

```
$> make debug
go run -mod vendor cmd/server/main.go \
		-verbose \
		-spatial-database-uri 'duckdb://?uri=/Users/asc/whosonfirst/go-whosonfirst-spatial-atproto/fixtures/sf_county.parquet'

2025/04/07 19:46:38 DEBUG Verbose logging enabled
2025/04/07 19:46:38 DEBUG Enable point in polygon handler endpoint=/xrpc/org.whosonfirst.PointInPolygon
2025/04/07 19:46:38 DEBUG Enable point in polygon from tile handler endpoint=/xrpc/org.whosonfirst.PointInPolygonWithTile
2025/04/07 19:46:38 DEBUG Enable point in polygon handler endpoint=/xrpc/org.whosonfirst.Intersects
2025/04/07 19:46:38 DEBUG Enable get record handler endpoint=/xrpc/org.whosonfirst.getRecord
2025/04/07 19:46:38 DEBUG Enable geocode handler endpoint=/xrpc/org.whosonfirst.geocode
2025/04/07 19:46:38 INFO Listening for requests address=http://localhost:8080
```

### Get record

Return a place by its unique ID.

```
$> curl 'http://localhost:8080/xrpc/org.whosonfirst.getRecord?id=102112179'
{
  "id": 102112179,
  "type": "Feature",
  "properties": {
    "geom:latitude": 37.748114,
    "geom:longitude": -122.420856,
    "wof:country": "US",
    "wof:id": 102112179,
    "wof:lastmodified": 1566604800,
    "wof:parent_id": 85922583,
    "wof:placetype": "neighbourhood",
    "wof:repo": "whosonfirst-data-admin-us"
  },
  "bbox": null,
  "geometry": {"coordinates":[[[-122.423623,37.739801],[-122.418225,37.748212],[-122.418106,37.748534],[-122.418419,37.751724],[-122.422614,37.7491],[-122.42221,37.745201],[-122.423305,37.742291],[-122.423885,37.74079],[-122.424104,37.739865],[-122.423623,37.739801]]],"type":"Polygon"}
}
```

#### Notes

As written this endpoint returns the raw GeoJSON record returned by the [whosonfirst/go-whosonfirst-spatial-duckdb](https://github.com/whosonfirst/go-whosonfirst-spatial-duckdb/blob/main/database_reader.go#L35) package. The use of GeoJSON in these responses is not to advocate for the format in ATProto/Geo responses but only to try and identify which properties a client may need to meet user-needs. For example, a "placetype" attribute to allow filtering for privacy or security reasons.

### Point in polygon

Return places that contain a given point.

```
$> curl -s 'http://localhost:8080/xrpc/org.whosonfirst.PointInPolygon?latitude=37.759991&longitude=-122.416977' | jq -r '.places[]["wof:name"]'

Inner Mission
Mission District
San Francisco
San Francisco
```

#### Notes

Should it be possible to filter (or exclude) results by placetype? Probably.

Under the hood this uses the [whosonfirst/go-whosonfirst-spatial/query.SpatialQuery](https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/query/query.go#L13) definition which defines additional query filters not exposed here. These include placetype, "current"-ness, dates, etc. It may be desirable to expose similar filtering criteria in ATProto/Geo queries.

As written this endpoint returns records encoded as a Who's On First [StandardPlacesResult](https://github.com/whosonfirst/go-whosonfirst-spr/blob/main/spr.go) (SPR). The goal behind the `SPR` was to define a minimum set of properties to be able to perform three functions:

1. Provide a minimum amount of data for filtering: placetype, is_current, etc.
2. Display a map with a point (centroid) and/or bounding box and a label.
3. Define URIs and endpoints where additional data may be retrieved.

The use of the `SPR` in these responses is not to advocate for the use of the Who's On First `SPR` in ATProto/Geo responses but only to try and identify which properties a client may need to meet user-needs. For example, a "placetype" attribute to allow filtering for privacy or security reasons.

### Intersects

Return places that intersect a given geometry.

```
$> curl -s 'http://localhost:8080/xrpc/org.whosonfirst.Intersects?geometry=POLYGON%28%28-122.423146+37.769809%2C-122.423301+37.771442%2C-122.423335+37.771478%2C-122.423238+37.771552%2C-122.423342+37.771697%2C-122.423397+37.771886%2C-122.42256+37.772528%2C-122.422547+37.771592%2C-122.421723+37.771627%2C-122.421635+37.770317%2C-122.421069+37.770129%2C-122.420227+37.769936%2C-122.420008+37.770111%2C-122.417753+37.769815%2C-122.415577+37.769591%2C-122.410931+37.769411%2C-122.410831+37.769232%2C-122.410814+37.769053%2C-122.408452+37.769163%2C-122.408174+37.769247%2C-122.408007+37.769244%2C-122.407853+37.768951%2C-122.407534+37.765783%2C-122.407305+37.763188%2C-122.40648+37.754564%2C-122.406453+37.754277%2C-122.406187+37.751367%2C-122.406224+37.751242%2C-122.406201+37.751064%2C-122.406164+37.750896%2C-122.40524+37.749128%2C-122.406312+37.748939%2C-122.406728+37.748819%2C-122.406758+37.748766%2C-122.407882+37.748405%2C-122.408706+37.748353%2C-122.413669+37.748255%2C-122.418225+37.748212%2C-122.418106+37.748534%2C-122.418419+37.751724%2C-122.422614+37.7491%2C-122.422879+37.751972%2C-122.423452+37.758367%2C-122.42805+37.758089%2C-122.428202+37.759692%2C-122.428354+37.76129%2C-122.430574+37.761155%2C-122.431025+37.765868%2C-122.430824+37.766014%2C-122.428949+37.767504%2C-122.428657+37.767735%2C-122.426309+37.769603%2C-122.423146+37.769809%29%29' | jq -r '.places[]["wof:name"]'

Lower Haight
The Castro
Duboce Triangle
Dolores Heights
Mission Dolores
Mint Hill
Baja Noe
Mission Dolores Park
Bernal Heights
La Lengua
Peralta Heights
Precitaville
Santana Rancho
Serpentinia
Inner Mission
Potrero Hill
Potrero Flats
Showplace Square
South of Market
West Soma
Mission District
San Francisco
San Francisco
```

#### Notes

Should this be a `GET` request? Probably not. Should we really be passing around URL-escaped [well-known text (WKT)](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry) geometries? Also, probably not. That said, I haven't read the ATProto docs enough to know what sort of restrictions and/or expectations there are around input parameters so for the sake of brevity (sort of) we'll just be "inelegant" about things.

Should it be possible to filter (or exclude) results by placetype? Probably.

Under the hood this uses the [whosonfirst/go-whosonfirst-spatial/query.SpatialQuery](https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/query/query.go#L13) definition which defines additional query filters not exposed here. These include placetype, "current"-ness, dates, etc. It may be desirable to expose similar filtering criteria in ATProto/Geo queries.

As written this endpoint returns records encoded as a Who's On First [StandardPlacesResult](https://github.com/whosonfirst/go-whosonfirst-spr/blob/main/spr.go) (SPR). The goal behind the `SPR` was to define a minimum set of properties to be able to perform three functions:

1. Provide a minimum amount of data for filtering: placetype, is_current, etc.
2. Display a map with a point (centroid) and/or bounding box and a label.
3. Define URIs and endpoints where additional data may be retrieved.

The use of the `SPR` in these responses is not to advocate for the use of the Who's On First `SPR` in ATProto/Geo responses but only to try and identify which properties a client may need to meet user-needs. For example, a "placetype" attribute to allow filtering for privacy or security reasons.

### Point in polygon with (map) tile

Return places that intersect a XZY map tile.

This method exists as a way to provide a privacy/security preserving means to return enough geographic data for user-defined extent (the bounding box of a map tile) rather than an exact coordinate such that a client may perform a final point-in-polygon operation on device.

_Note: As of this writing the implementation of this method does NOT return the necessary data to perform an on-device point-in-polygon operation. The current implementation is mostly to imagine what the chatter between a client and (geo) service provider might look like._

```
$> curl -s 'http://localhost:8080/xrpc/org.whosonfirst.PointInPolygonWithTile?z=12&x=655&y=1583' | jq -r '.places[]["wof:name"]'                                      
Excelsior
Mission Terrace
St. Mary's Park
St. Mary's
Lost Tribe of College Hill
Little Saigon
The Sit/Lie
Off Market
Le marché
Little Saigon
Union Square
Forgotten Island
The Park
The Hydeaway
The Rambles
Tender Wasteland
The Loin Pit
The Panhandle
The Naked Hood
Tenderloin East
Delicious Fields
The Gimlet
French Quarter
Glen Park
Western Addition
Noe Valley
Alamo Square
Lower Haight
The Castro
Duboce Triangle
Japantown
Dolores Heights
Hayes Valley
Mission Dolores
Mint Hill
Cathedral Hill
Fairmount
Baja Noe
Thomas Paine Square
Little Osaka
Laguna Heights
St. Francis Square
Malcolm X Square
Mission Dolores Park
Lower Nob Hill
Rental Row
The Whoa-Man
The Post Up
Tenderloin Heights
Academy Downs
BoHo Slope
Visitacion Valley
Bayview Heights
Hunters Point
Portola
Bayview
Silver Terrace
India Basin
Bernal Heights
Sutro Vista
Hill People of Powhattan
Cortlandia
Alemanistan
Baja Cortlandia
Holly Park
The Crescent
Liminal Zone of Deliciousness
University Mound
Apparel City
Bayview-Hunters Point
South Basin
Candlestick Point
Bret Harte
La Lengua
Eugeniaia
Sutrito Canine Republic
Esmereldia
Peralta Heights
NanoTokyo
Precitaville
Santana Rancho
Serpentinia
Principality of Chicken John
Inner Mission
Potrero Hill
Potrero Flats
Showplace Square
South of Market
Intermission
Civic Center
Financial District South
The Nipple
Foxy Heights
Civic Center
Opera Plaza
The Money Shot
Ghost Market
Fecal Fountain
The Yo
Tenderloin
Saint Anne's
Van Ness
The Bar
The Castle Triangle
Pill Hill
Deli Hills
Tender Turnpike
NFL Experience
Victoria Mews
West Soma
Financial District
Mission Bay
Central Waterfront
Rincon Hill
South Beach
Dogpatch
McLaren Park
Downtown
Mission District
Produce Market
San Francisco
San Francisco
```

#### Notes

Should it be possible to filter (or exclude) results by placetype? Probably.

As written this endpoint returns records encoded as a Who's On First [StandardPlacesResult](https://github.com/whosonfirst/go-whosonfirst-spr/blob/main/spr.go) (SPR). The goal behind the `SPR` was to define a minimum set of properties to be able to perform three functions:

1. Provide a minimum amount of data for filtering: placetype, is_current, etc.
2. Display a map with a point (centroid) and/or bounding box and a label.
3. Define URIs and endpoints where additional data may be retrieved.

The use of the `SPR` in these responses is not to advocate for the use of the Who's On First `SPR` in ATProto/Geo responses but only to try and identify which properties a client may need to meet user-needs. For example, a "placetype" attribute to allow filtering for privacy or security reasons.

As mentioned earlier the `SPR` is not a good fit for this operation since it only returns bounding boxes and not actually geometries necessary to perform a point-in-polygon operation on device.

### Geocode

Returns places matching a string (geocode).

```
$> curl -s 'http://localhost:8080/xrpc/org.whosonfirst.geocode?q=Latin+American+Club' | jq
[
  {
    "id": 571986789,
    "name": "Latin American Club",
    "placetype": "venue",
    "lineage": [
      {
        "locality": {
          "id": 85922583,
          "name": "San Francisco",
          "abbr": "SF",
          "languageDefaulted": true
        },
        "macrohood": {
          "id": 1108830809,
          "name": "Mission District",
          "languageDefaulted": true
        },
        "neighbourhood": {
          "id": 85834637,
          "name": "Inner Mission",
          "languageDefaulted": true
        },
        "venue": {
          "id": 571986789,
          "name": "Latin American Club",
          "languageDefaulted": true
        }
      }
    ],
    "geom": {
      "bbox": "-122.420536,37.755348,-122.420536,37.755348",
      "lat": 37.755348,
      "lon": -122.420536
    },
    "languageDefaulted": true
  }
]
```

#### Notes

This method simply proxies the [Placeholder API](https://github.com/pelias/placeholder/). As such you will need to run an instance of the Placeholder server and specify its endpoint with the `-placeholder-endpoint` command-line flag (both default to `http://localhost:3000`).

The Placeholder API uses a SQLite database for lookups. The default SQLite database is 5GB with global coverage. For the purposes of this demo there is a smaller (San Francisco country sized) database that is included in the `fixtures` folder. This database contains both administrative records _and_ venues. Like the GeoParquet files this database is bundled using Git LFS and is further compressed using `bzip2`.

To use the bundled Placeholder database you will do to the following (adjusting paths as necessary):

```
$> bunzip -k fixtures/store.sqlite3.bz2

$> cd /usr/local/placeholder
$> npm install
$> export PLACEHOLDER_DATA /usr/local/go-whosonfirst-spatial-atproto/fixtures
$> export HOST localhost
$> npm start
```

As written this endpoint returns records encoded as Placeholder API results. Their use is not to advocate for the use of the Who's On First `SPR` in ATProto/Geo responses but only to try and identify which properties a client may need to meet user-needs. For example, a "placetype" attribute to allow filtering for privacy or security reasons.

## See also

* https://github.com/schuyler/garganorn
* https://atproto.com/specs/xrpc
* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-www
* https://github.com/whosonfirst/go-whosonfirst-spatial-duckdb
