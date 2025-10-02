package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andythigpen/clock2/pkg/services"
)

type DisplayStateHandler struct {
	svc *services.DisplayService
}

func NewDisplayStateHandler(svc *services.DisplayService) *DisplayStateHandler {
	return &DisplayStateHandler{svc}
}

func (h *DisplayStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		state, err := h.svc.GetState()
		if err != nil {
			slog.Error("failed to get display state", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(map[string]string{"state": string(state)}); err != nil {
			slog.Error("failed to encode state", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		dec := json.NewDecoder(r.Body)
		body := map[string]string{}
		if err := dec.Decode(&body); err != nil {
			slog.Error("failed to decode body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var state services.DisplayState
		switch strings.ToLower(body["state"]) {
		case "on":
			state = services.DisplayStateOn
		case "off":
			state = services.DisplayStateOff
		default:
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("invalid state in request", "body", body)
			return
		}
		if err := h.svc.SetState(state); err != nil {
			slog.Error("failed to set display state", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

type DisplayBrightnessHandler struct {
	svc        *services.DisplayService
	brightness uint8
	updateCh   chan uint8
}

func NewDisplayBrightnessHandler(svc *services.DisplayService) *DisplayBrightnessHandler {
	return &DisplayBrightnessHandler{svc: svc, brightness: 100, updateCh: make(chan uint8)}
}

func (h *DisplayBrightnessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := r.URL.Query()
		current := q.Get("current")
		if current != "" {
			b, err := strconv.ParseInt(current, 10, 8)
			if err != nil {
				slog.Error("invalid brightness parameter", "current", current)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// http long polling
			if b == int64(h.brightness) {
				select {
				case b := <-h.updateCh:
					slog.Debug("received brightness notification", "b", b)
				case <-time.After(30 * time.Second):
					slog.Debug("no brightness notification received")
				}
			}
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(map[string]uint8{"brightness": h.brightness}); err != nil {
			slog.Error("failed to encode brightness", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		dec := json.NewDecoder(r.Body)
		body := map[string]any{}
		if err := dec.Decode(&body); err != nil {
			slog.Error("failed to decode body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if b, ok := body["brightness"]; ok {
			f, ok := b.(float64) // the json encoder treats numbers as floats
			if !ok {
				slog.Error("expected an number for brightness", "body", body)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			i := uint8(f)
			if i > 100 {
				slog.Error("invalid brightness", "brightness", i)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			h.brightness = i
			h.svc.SetBrightness(i)
			select {
			case h.updateCh <- h.brightness:
				slog.Debug("notify brightness change", "brightness", h.brightness)
			default:
				slog.Debug("no brightness change receivers", "brightness", h.brightness)
			}
		}
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
