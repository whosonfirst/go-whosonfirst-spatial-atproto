# ATGeo lexicon

_This is an incomplete draft. If you're reading this understand that it is work-in-progress and incomplete._

```
type Status uint8

const ( _ Status = iota
      Current
      Retired
      Superseded
      Deprecated
)

type Place struct {
     ID string `json:"id"`
     URI string `json:"uri"`
     Name string `json:"name"`
     Placetype string `json:"placetype"`
}

type Location struct {
     Place
     Status Status `json:"status"`     
     Hierarchy []*Place `json:"hierarchy"`
     SupersededBy []string `json:"superseded_by"`
     Supersedes []string `json:"supersedes"`
     Geometry *geojson.Geometry `json:"geometry,omitempty"`
}
```

## Place

A "place" is a minimum set of properties shared between both a location record and other records that parent or are ancestors of that location. That is to say an ancestor referenced in a location record has _less_ information that the location itself.

### ID

A stable, permanent (canonical) identifier for this place within the context of the location provider (gazetteer).

### URI

The HTTP-addressable URI where this record can be resolved to in full. That fully-resolved record may, or may not, be an ATGeo lexicon response.

#### Notes

In this example, the semantics of the ID "provider" are assumed to be accounted for by the TLD of the domain. That might need to be revisited.

### Name

The principal name for this location record. This is distinct from a more complete label. For example "Montréal" (name) versus "Montréal, QC, Canada".

#### Notes

But what language? Exactly. The mechanics of specifying one or more names remains to be worked out.

The approach that the Who's On First project has taken has been to say that every record has a `wof:name` property in the "default" language and then zero or more language/dialect specific `name:` properties. For example: [San Francico](https://spelunker.whosonfirst.org/id/85922583/geojson).

Importantly as of this writing the "default" language is English but in the future it may be something else. The point is to enforce a common label across all records.

### Placetype


## Location

These are properties specific to a location associated with an ATProto message/event/whatever.

### Status

### SupersededBy

### Supersedes

### Latitude

### Geometry

