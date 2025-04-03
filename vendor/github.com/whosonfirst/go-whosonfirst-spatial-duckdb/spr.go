package duckdb

import (
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

// This exists only because there isn't a default implementation of spr.StandardPlacesResults
// provided by whosonfirst/go-whosonfirst-spr/v2 and there should be.

type SPRResults struct {
	spr.StandardPlacesResults `json:",omitempty"`
	SPRResults                []spr.StandardPlacesResult `json:"places"`
}

func (r *SPRResults) Results() []spr.StandardPlacesResult {
	return r.SPRResults
}

func NewSPRResults(results []spr.StandardPlacesResult) spr.StandardPlacesResults {

	return &SPRResults{
		SPRResults: results,
	}
}
