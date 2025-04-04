package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/whosonfirst/go-whosonfirst-spatial-atproto/xrpc"
)

func xrpcError(rsp http.ResponseWriter, message string, code int) {

	xrpc_err := xrpc.Error{
		Error:   strconv.Itoa(code),
		Message: message,
	}

	rsp.Header().Set("Content-type", "application/json")

	enc := json.NewEncoder(rsp)
	err := enc.Encode(xrpc_err)

	if err != nil {
		slog.Error("Failed to encode XRPC error", "error", err)
		http.Error(rsp, "Internal server error", http.StatusInternalServerError)
	}

	return
}
