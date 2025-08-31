package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/andythigpen/clock2/pkg/handlers"
	"github.com/andythigpen/clock2/pkg/services"
)

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

	ctx := context.Background()

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

	mux := http.NewServeMux()
	handlers.Register(mux, haSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	slog.Info("listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("failed to listen", "addr", addr)
	}
	haSvc.Stop()
}
