package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/andythigpen/clock2/pkg/handlers"
	"github.com/andythigpen/clock2/pkg/services"
	"github.com/andythigpen/clock2/pkg/ui"
)

func makeHomeAssistant(ctx context.Context) *services.HomeAssistantService {
	haUrl := os.Getenv("HA_URL")
	haToken := os.Getenv("HA_TOKEN")
	weatherEntity := os.Getenv("HA_WEATHER_ENTITY")
	forecastEntity := os.Getenv("HA_FORECAST_ENTITY")
	sunEntity := os.Getenv("HA_SUN_ENTITY")
	haSvc := services.NewHomeAssistantService(
		haUrl,
		haToken,
		services.WithWeatherEntity(weatherEntity),
		services.WithForecastEntity(forecastEntity),
		services.WithSunEntity(sunEntity),
	)
	go haSvc.RunForever(ctx)
	return haSvc
}

func makeDisplayService() *services.DisplayService {
	return services.NewDisplayService()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "warning", "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error", "err":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	flag.Parse()

	ctx := context.Background()
	haSvc := makeHomeAssistant(ctx)
	displaySvc := makeDisplayService()

	mux := http.NewServeMux()
	handlers.Register(mux, displaySvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	slog.Info("listening", "addr", addr)
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to listen", "addr", addr)
		}
	}()

	ui.RunForever(haSvc, displaySvc)

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server")
	}
	haSvc.Stop()
}
