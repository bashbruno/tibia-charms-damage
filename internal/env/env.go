package env

import (
	"log/slog"
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("Failed to read environment variable, using fallback", "key", key, "fallback", fallback)
		return fallback
	}
	return v
}

func GetInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("Failed to read environment variable, using fallback", "key", key, "fallback", fallback)
		return fallback
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		slog.Warn("Failed to parse environment variable as integer, using fallback", "key", key, "value", v, "fallback", fallback, "error", err)
		return fallback
	}

	return i
}
