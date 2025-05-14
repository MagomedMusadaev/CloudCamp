package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"CloudCamp/internal/config"
)

// SetupLogger настраивает глобальный логгер через slog.SetDefault
func SetupLogger(cfg *config.Config) (io.Closer, error) {
	const op = "internal.logger.SetupLogger"

	// Создаём директорию для логов, если её нет
	if err := os.MkdirAll(cfg.Log.GetLogDir(), 0755); err != nil {
		return nil, fmt.Errorf("%s: failed to create log dir: %w", op, err)
	}

	// Открываем (или создаём) лог-файл
	logFile, err := os.OpenFile(
		cfg.Log.FilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open log file: %w", op, err)
	}

	var handler slog.Handler

	switch cfg.Env {
	case config.EnvDev:
		handler = slog.NewTextHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{
			Level: slog.LevelDebug,
			//AddSource: true,
		})
	case config.EnvProd:
		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case config.EnvTest:
		handler = slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelError,
		})
	default:
		return nil, fmt.Errorf("%s: unknown environment %q", op, cfg.Env)
	}

	// Создаём и устанавливаем глобальный логгер
	slog.SetDefault(slog.New(handler))

	return logFile, nil
}
