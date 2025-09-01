package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/andythigpen/clock2/pkg/components"
	"github.com/andythigpen/clock2/pkg/services"
)

func Register(mux *http.ServeMux, haSvc *services.HomeAssistantService, displaySvc *services.DisplayService) {
	component := components.Index()
	mux.Handle("/", templ.Handler(component))
	mux.Handle("/carousel", NewCarouselHandler(haSvc, displaySvc))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/api/display/state", NewDisplayStateHandler(displaySvc))
	mux.Handle("/api/display/brightness", NewDisplayBrightnessHandler())
}
