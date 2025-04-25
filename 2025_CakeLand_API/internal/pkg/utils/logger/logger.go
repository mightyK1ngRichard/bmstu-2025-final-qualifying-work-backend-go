package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type EnvKind string

const (
	Local EnvKind = "local"
	Dev   EnvKind = "dev"
	Prod  EnvKind = "prod"
)

func NewLogger(envKind EnvKind) *slog.Logger {
	var log *slog.Logger

	switch envKind {
	case Local:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case Dev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case Prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func (e *EnvKind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}

	switch EnvKind(value) {
	case Local, Dev, Prod:
		*e = EnvKind(value)
		return nil
	default:
		return fmt.Errorf("неизвестное значение окружения: %s", value)
	}
}
