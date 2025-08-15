package main

import (
	"log/slog"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	slog.Error("Err", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "error", err.Error())
}
