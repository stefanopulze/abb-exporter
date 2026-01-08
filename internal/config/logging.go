package config

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Type  string `yaml:"type" env:"LOG_TYPE" env-default:"console"`
}

func configureLogging(config LogConfig) {
	var handler slog.Handler
	logConfig := &slog.HandlerOptions{
		Level: logLevel(config.Level),
	}

	if config.Type == "console" {
		handler = slog.NewTextHandler(os.Stdout, logConfig)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, logConfig)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func logLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
