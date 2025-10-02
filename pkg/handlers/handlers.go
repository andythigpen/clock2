package handlers

import (
	"net/http"

	"github.com/andythigpen/clock2/pkg/services"
)

func Register(
	mux *http.ServeMux,
	displaySvc *services.DisplayService,
) {
	mux.Handle("/api/display/state", NewDisplayStateHandler(displaySvc))
	mux.Handle("/api/display/brightness", NewDisplayBrightnessHandler(displaySvc))
}
