package main

import (
	"CloudCamp/internal/app"
	"CloudCamp/internal/background"
	"CloudCamp/internal/config"
	"CloudCamp/internal/logger"
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Парсим флаги командной строки для получения пути конфиг.yaml файла
	configPath := flag.String("config", "configs/config.yaml", "path to configuration file")
	flag.Parse()

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	// Настройка логгера
	closed, err := logger.SetupLogger(cfg)
	if err != nil {
		slog.Error("Error setting up logger", "error", err)
		os.Exit(1)
	}
	defer closed.Close()

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем сервер
	server, err := app.NewServer(cfg)
	if err != nil {
		slog.Error("Error creating server", "error", err)
		os.Exit(1)
	}

	// Создаем и запускаем фоновые процессы
	healthChecker := background.NewHealthChecker(
		server.GetBackends(),
		cfg.HealthChecker.Interval,
		cfg.HealthChecker.Path,
	)

	healthChecker.Start(ctx)

	// Канал для получения сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err = server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error running server", "error", err)
			cancel()
		} else {
			slog.Info("Server shutdown gracefully")
		}
	}()

	// Ожидаем сигнала завершения
	<-sigChan
	slog.Info("Shutting down...")

	// Отменяем контекст для остановки фоновых процессов
	cancel()

	// Ожидаем завершения фоновых процессов
	healthChecker.Wait()

	// Останавливаем сервер
	if err = server.Shutdown(); err != nil {
		slog.Error("Error shutting down server", "error", err)
	}
}
