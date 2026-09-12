package engine

import (
	"log/slog"
	"os"
	"time"

	"github.com/AliAlbhrani/StudentsArchive/env"
	"github.com/lmittmann/tint"
)

var _ = func() bool {
	var handler slog.Handler
	if env.DEV {
		handler = tint.NewTextHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		})
	}

	slog.SetDefault(slog.New(handler))
	return true
}()
