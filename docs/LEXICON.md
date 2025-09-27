# ATGeo lexicon

## Motivation

This is an initial draft proposal for an ATGeo location record (lexicon). In fact, it's more of a suggestion rather than a proposal for the purposes or stimulating discussion and understanding issues and requirements.

It is a simplified version of the Who's On First "standard places response".

## Model

Individual properties are discussed in detail below.

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
     Hierarchy []*Place `json:"hierarchy"`     
     Status Status `json:"status"`     
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

`ID` and `URI` could of course be simplified to a single property: A namespace prefix, mapped to a URI, and an identifier separated by a colon but then you either have to choose between all the hassle and complexity of XML namespaces or the ease, but potential ambiguity, of Flickr-style machine tags.

Who's On First addresses this issue by endeavouring to map all properties and their sources (namespaces) to machine-readable documents that can be derived from reliable URI templates. For example `wof:placetype` maps to:

```
https://github.com/whosonfirst/whosonfirst-properties/blob/main/properties/wof/placetype.json
``

And `wof` maps to:

```
https://github.com/whosonfirst/whosonfirst-sources/blob/main/sources/wof.json
```

### Name

The principal name for this location record. This is distinct from a more complete label. For example "Montréal" (name) versus "Montréal, QC, Canada".

#### Notes

But what language? Exactly. The mechanics of specifying one or more names remains to be worked out.

The approach that the Who's On First project has taken has been to say that every record has a `wof:name` property in the "default" language and then zero or more language/dialect specific `name:` properties. For example: [San Francico](https://spelunker.whosonfirst.org/id/85922583/geojson).

Importantly as of this writing the "default" language is English but in the future it may be something else. The point is to enforce a common label, with consistent semantics, across all records.

Assuming that an ATGeo record does _not_ want to enforce a "default" language then it stands to reason that `name` should not be a string but rather a struct containing label and language details.

As with the status property, discussed below, then what all of this suggests is (possibly) the need for language-specific name/label xRPC lookup methods.

As mentioned the `name` property is not meant to be an application-specific label nor is it meant to encode location-specific metadata, for example the address of a venue.

### Placetype

The type, or descriptor, for this location.

#### Note

Like names, placetypes are harder than anyone would like. In this example "placetype" is defined as a string, as opposed to fixed list, which opens it up to being a free-for-all.

Who's On First takes a different approach. In the WOF model there are three different types of "places": Common, optional and common optional. Any place can have any one of those placetypes and can have as complex a hierarchy of ancestor (placetypes) as necessary to represent its reality.

The only rule is that every place has a _minimum_ of one _common_ placetype in its hierarchy. This ensures that any two (or more) projects have a shared set of ("common") placetypes that they can use to match place records regardless of the details or nuanced required by anyone project.

Who's On First placetypes are discussed in detail here: https://github.com/whosonfirst/whosonfirst-placetypes

Like sources and properties, every place type has a machine-readable representation. For example:

https://github.com/whosonfirst/whosonfirst-placetypes/blob/main/placetypes/campus.json

Although these records do not currently have language-specific `name:` translations there's nothing preventing that happening.

## Location

These are properties specific to a location associated with an ATProto message/event/whatever.

### Hierarchy

An ordered list of ancestors that a location is "parented" by. This is meant to be used to construct application-specific labels for a location. For example "Vancouver, British Columbia CA" rather than "Vancouver".

#### Notes

The property has all the same language-related issues of the `name` property.

### Status

An enumerated list of possible states for a location record: current, retired, superseded, deprecated.

#### Notes

I am not convinced this is really necessary or practical in an ATGeo record since once embedded in any given ATProto message there's no way to update it after the fact (for example if a given location record is deprecated or otherwise retired).

This is perhaps better suited to a geo-specific xRPC method associated with location records?

### SupersededBy

One or more location records that supersede this location.

#### Notes

As with the status property it's not clear to me that this makes much sense in an ATgeo record.

### Supersedes

One or more location records that this location supersedes.

#### Notes

As with the status property it's not clear to me that this makes much sense in an ATgeo record.

### Geometry

A GeoJSON `geometry` element. Importantly, for privacy reasons, this property is _OPTIONAL_. Importantly this means that ATGeo location records are associated with the "idea" of a place, as defined by its `id` and `uri` properties, rather than its geographic representation.

#### Notes

Allowing for any valid A GeoJSON `geometry` element allows records to encode complex geographic data which better represents its reality. In the case of administrative locations this also allows for a minimum-bounding rectangle to be derived which allows applications to more accurately visualize that location on a map.

At the same time the absence of explicit "centroid" coordinates may prove problematic for applications. Simply deriving the "center" of a complex geometry is not always the best, or correct, point on which a human-readable label should be placed on a map. For example the geographic center of the city of San Francisco is 15 miles west of the city, in a Pacific ocean, since the city, as a legal entity, also encompasses the Farralon Islands located 30 miles of the coast.

Who's On First addresses this issue by overloading the term "centroid" to mean "area of focus" and then defines one or more of a series of named-centroids to associate with any given record. Those named centroids are descibed here:

https://github.com/whosonfirst/whosonfirst-geometries/blob/main/geometries/README.md


