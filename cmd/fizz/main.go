package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/laouji/fizz/pkg/cache"
	"github.com/laouji/fizz/pkg/handler"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type config struct {
	WebPort       int
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

func main() {
	logger := logrus.New()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	cfg, err := loadConfig()
	if err != nil {
		logger.WithError(err).Errorf("failed to load environment variables")
		os.Exit(1)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0, // use default DB
	})

	server, err := configureServer(cfg.WebPort, redisClient, logger)
	if err != nil {
		logger.WithError(err).Errorf("failed configure server")
		os.Exit(2)
	}

	go func() {
		logger.Info("listening...")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Errorf("failed to listen to the given port: %d", cfg.WebPort)
			os.Exit(2)
		}
		logger.Info("web server stopped receiving new conns")
	}()

	// block until signal is received
	<-signalCh

	shutdownCtx, release := context.WithTimeout(context.Background(), 3*time.Second)
	defer release()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("failed to gracefully shutdown server")
		os.Exit(2)
	}
	logger.Info("graceful shutdown complete")
}

func loadConfig() (*config, error) {
	portStr := os.Getenv("WEB_PORT")
	webPort, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	portStr = os.Getenv("REDIS_PORT")
	redisPort, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &config{
		WebPort:       webPort,
		RedisPort:     redisPort,
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}, nil
}

func configureServer(
	port int,
	redisClient *redis.Client,
	logger logrus.FieldLogger,
) (*http.Server, error) {
	logger.Info("initializing routes")

	cache := cache.NewClient(redisClient, logger)

	r := mux.NewRouter()
	r.HandleFunc("/stats", handler.Stats(cache, logger))
	r.HandleFunc("/", handler.FizzBuzz(cache, logger))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	return server, nil
}
