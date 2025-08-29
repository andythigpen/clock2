package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"

	"github.com/andythigpen/clock2/pkg/components"
)

func main() {
	component := components.Index()

	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(component))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	slog.Info("Listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("Failed to listen", "addr", addr)
	}
}
