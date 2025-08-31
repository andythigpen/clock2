package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/andythigpen/clock2/pkg/components"
	"github.com/andythigpen/clock2/pkg/services"
)

func Register(mux *http.ServeMux, haSvc *services.HomeAssistantService) {
	component := components.Index()
	mux.Handle("/", templ.Handler(component))
	mux.Handle("/carousel", NewCarouselHandler(haSvc))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}
