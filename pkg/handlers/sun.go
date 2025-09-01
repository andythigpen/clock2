package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andythigpen/clock2/pkg/services"
)

type SunHandler struct {
	svc *services.HomeAssistantService
}

func NewSunHandler(svc *services.HomeAssistantService) *SunHandler {
	return &SunHandler{svc}
}

func (h *SunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	sun := h.svc.GetSun()
	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]any{
		"nextRising":  sun.Attributes.NextRising,
		"nextSetting": sun.Attributes.NextSetting,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
