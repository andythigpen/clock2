package homeassistant

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/andythigpen/clock2/pkg/models/weather"
)

type entityState[T any] struct {
	id    string
	state T
}

type homeAssistantService struct {
	baseUrl  string
	token    string
	doneCh   chan bool
	weather  entityState[weather.WeatherEntity]
	forecast entityState[weather.ForecastEntity]
	sun      entityState[weather.SunEntity]
}

type HomeAssistantServiceOption func(*homeAssistantService)

func WithWeatherEntity(id string) HomeAssistantServiceOption {
	return func(r *homeAssistantService) {
		r.weather.id = id
	}
}

func WithForecastEntity(id string) HomeAssistantServiceOption {
	return func(r *homeAssistantService) {
		r.forecast.id = id
	}
}

func WithSunEntity(id string) HomeAssistantServiceOption {
	return func(r *homeAssistantService) {
		r.sun.id = id
	}
}

func NewHomeAssistantService(baseUrl, token string, opts ...HomeAssistantServiceOption) homeAssistantService {
	svc := homeAssistantService{
		baseUrl: baseUrl,
		token:   token,
		doneCh:  make(chan bool),
	}
	for _, o := range opts {
		o(&svc)
	}
	return svc
}

func (r *homeAssistantService) RunForever(ctx context.Context) {
	slog.Info("running HomeAssistant service")
	// fetch the initial state
	r.poll(ctx)
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-r.doneCh:
			ticker.Stop()
			return
		case <-ticker.C:
			r.poll(ctx)
		}
	}
}

func (r *homeAssistantService) Stop() {
	r.doneCh <- true
}

func (r *homeAssistantService) GetWeather() weather.WeatherEntity {
	return r.weather.state
}

func (r *homeAssistantService) GetForecast() weather.ForecastEntity {
	return r.forecast.state
}

func (r *homeAssistantService) GetSun() weather.SunEntity {
	return r.sun.state
}

func (r *homeAssistantService) poll(ctx context.Context) {
	if r.weather.id != "" {
		if err := r.fetchEntity(ctx, r.weather.id, &r.weather.state); err != nil {
			slog.ErrorContext(ctx, "failed to fetch entity", "entityId", r.weather.id, "err", err)
		}
	}

	if r.forecast.id != "" {
		if err := r.fetchEntity(ctx, r.forecast.id, &r.forecast.state); err != nil {
			slog.ErrorContext(ctx, "failed to fetch entity", "entityId", r.forecast.id, "err", err)
		}
	}

	if r.sun.id != "" {
		if err := r.fetchEntity(ctx, r.sun.id, &r.sun.state); err != nil {
			slog.ErrorContext(ctx, "failed to fetch entity", "entityId", r.sun.id, "err", err)
		}
	}
}

func (r *homeAssistantService) fetchEntity(ctx context.Context, entityId string, v any) error {
	u, err := url.JoinPath(r.baseUrl, "/api/states", entityId)
	if err != nil {
		return err
	}
	slog.Debug("fetching entity", "entityId", entityId, "url", u)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
