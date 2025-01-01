package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/TaylorOno/wg/drinks"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevelFromEnv()}))
	slog.SetDefault(logger)

	listenAddr := ":9090"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/wg-drinks/menu", drinks.DrinkMenu)
	http.HandleFunc("POST /wg-drinks/history", drinks.History)
	http.HandleFunc("GET /wg-drinks/history/{id}", drinks.History)
	slog.Info(fmt.Sprintf("listening on %s. Go to http://127.0.0.1%s/", listenAddr, listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func getLogLevelFromEnv() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo

	}
}
