package api

import (
	// "encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http/v3/sanitize"
	"github.com/whosonfirst/go-whosonfirst-reader/v2"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
)

type GetRecordHandlerOptions struct{}

func GetRecordHandler(app *spatial_app.SpatialApplication, opts *GetRecordHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()
		ctx := req.Context()

		id, err := sanitize.GetInt64(req, P_RECORD_ID)

		if err != nil {
			logger.Error("Failed to derive record ID", "error", err)
			xrpcError(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		body, err := reader.LoadBytes(ctx, app.SpatialDatabase, id)

		if err != nil {
			logger.Error("Failed to load record", "id", id, "error", err)
			xrpcError(rsp, "Not found", http.StatusNotFound)
			return
		}

		// Generate ATGeo here...

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write(body)
		return
	}

	record_handler := http.HandlerFunc(fn)
	return record_handler, nil
}
