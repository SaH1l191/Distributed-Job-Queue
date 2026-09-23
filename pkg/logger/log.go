package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func InitLogger() *slog.Logger {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "03:04PM",
		AddSource:  false,
	})
	return slog.New(handler)
}
