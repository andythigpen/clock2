package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"

	"github.com/andythigpen/clock2/pkg/handlers"
	"github.com/andythigpen/clock2/pkg/services"
)

//go:embed assets/*
var assets embed.FS

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
	opts := []services.DisplayServiceOption{}
	displayCmd := os.Getenv("DISPLAY_CMD")
	if displayCmd != "" {
		opts = append(opts, services.WithDisplayCommand(displayCmd))
	}
	getArgs := os.Getenv("DISPLAY_GET_ARGS")
	if getArgs != "" {
		opts = append(opts, services.WithDisplayGetArgs(strings.Split(getArgs, " ")...))
	}
	onArgs := os.Getenv("DISPLAY_ON_ARGS")
	if onArgs != "" {
		opts = append(opts, services.WithDisplayOnArgs(strings.Split(onArgs, " ")...))
	}
	offArgs := os.Getenv("DISPLAY_OFF_ARGS")
	if offArgs != "" {
		opts = append(opts, services.WithDisplayOffArgs(strings.Split(offArgs, " ")...))
	}
	onMatch := os.Getenv("DISPLAY_ON_MATCH")
	if onMatch != "" {
		opts = append(opts, services.WithDisplayOnMatch(regexp.MustCompile(onMatch)))
	}
	offMatch := os.Getenv("DISPLAY_OFF_MATCH")
	if offMatch != "" {
		opts = append(opts, services.WithDisplayOffMatch(regexp.MustCompile(offMatch)))
	}
	return services.NewDisplayService(opts...)
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

	ctx := context.Background()
	haSvc := makeHomeAssistant(ctx)
	displaySvc := makeDisplayService()

	mux := http.NewServeMux()
	handlers.Register(mux, haSvc, displaySvc, assets)

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
