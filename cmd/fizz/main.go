package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/laouji/fizz/pkg/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	port := 5000
	server, err := configureServer(port, logger)
	if err != nil {
		logger.WithError(err).Errorf("failed configure server")
		os.Exit(2)
	}

	go func() {
		logger.Info("listening...")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Errorf("failed to listen to the given port: %d", port)
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

func configureServer(
	port int,
	logger logrus.FieldLogger,
) (*http.Server, error) {
	logger.Info("initializing routes")
	r := mux.NewRouter()
	r.HandleFunc("/", handler.FizzBuzz(logger))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	return server, nil
}
