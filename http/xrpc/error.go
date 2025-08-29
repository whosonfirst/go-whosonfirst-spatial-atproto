package xrpc

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type xrpcError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Error(rsp http.ResponseWriter, message string, code int) {

	xrpc_err := xrpcError{
		Error:   strconv.Itoa(code),
		Message: message,
	}

	rsp.Header().Set("Content-type", "application/json")

	// Unclear if I am supposed to assign an HTTP error here
	// or if they've gone full-Flickr and are returning
	// errors in HTTP 200 OK responses...

	enc := json.NewEncoder(rsp)
	err := enc.Encode(xrpc_err)

	if err != nil {
		slog.Error("Failed to encode XRPC error", "error", err)
		http.Error(rsp, "Internal server error", http.StatusInternalServerError)
	}

	return
}
