package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dsha256/packer/internal/packer"
)

const (
	version = "1.0.0"

	idleTimeout  = 1 * time.Minute
	readTimeout  = 10 * time.Second
	writeTimeout = 30 * time.Second

	gracefulShutdownTimeout = 5 * time.Second
)

// Server holds params for REST API server configuration.
type Server struct {
	SizerSrvc  *packer.SizerService
	PackerSrvc *packer.PacketsService
}

// NewServer constructs Server instance.
func NewServer(sizerSrvc *packer.SizerService, packerSrvc *packer.PacketsService) *Server {
	return &Server{
		SizerSrvc:  sizerSrvc,
		PackerSrvc: packerSrvc,
	}
}

// Serve runs REST API server.
func (s *Server) Serve(port string) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      s.routes(),
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		qs := <-quit
		slog.Info("shutting down server",
			slog.Any("signal", qs.String()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)
	}()

	slog.Info("starting server",
		slog.Any("addr", srv.Addr),
	)

	// Calling Shutdown() on the server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("stopped server",
		slog.Any("addr", srv.Addr),
	)

	return nil
}
