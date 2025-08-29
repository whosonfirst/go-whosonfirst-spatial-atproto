package whosonfirst

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/aaronland/go-http/v3/sanitize"
	"github.com/whosonfirst/go-whosonfirst-spatial-atproto/http/xrpc"
)

type GeocodeHandlerOptions struct {
	PlaceholderEndpoint string
}

func GeocodeHandler(opts *GeocodeHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		query, err := sanitize.GetString(req, P_QUERY)

		if err != nil {
			logger.Error("Failed to derive query", "error", err)
			xrpc.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		u, err := url.Parse(opts.PlaceholderEndpoint)

		if err != nil {
			logger.Error("Failed to parse placeholder endpoint", "error", err)
			xrpc.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		u.Path = "/parser/search"

		q := url.Values{}
		q.Set("text", query)

		u.RawQuery = q.Encode()

		ph_rsp, err := http.Get(u.String())

		if err != nil {
			logger.Error("Failed to query placeholder endpoint", "error", err)
			xrpc.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer ph_rsp.Body.Close()

		if ph_rsp.StatusCode != http.StatusOK {
			logger.Error("Placeholder did not return ok", "code", ph_rsp.StatusCode)
			xrpc.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Generate ATGeo here...

		rsp.Header().Set("Content-type", "application/json")

		_, err = io.Copy(rsp, ph_rsp.Body)

		if err != nil {
			xrpc.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	geocode_handler := http.HandlerFunc(fn)
	return geocode_handler, nil
}
