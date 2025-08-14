package env

import (
	"log/slog"
	"os"
)

func GetString(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		slog.Warn("Failed to read environment variable, using fallback", "key", key, "fallback", fallback)
		return fallback
	}
	return v
}
