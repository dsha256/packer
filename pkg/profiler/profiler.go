package profiler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"time"
)

var ErrServerAlreadyStarted = errors.New("profiler server already started")

const (
	defaultHTTPPort                 = 4667
	defaultHTTPReadHeaderTimeout    = 3 * time.Second
	defaultURLPathFirstSubdirectory = "pprof"
)

// Config holds the configuration for the profiler.
type Config struct {
	HTTPServer               *http.Server
	URLPathFirstSubdirectory string
	HTTPPort                 int
	HTTPReadHeaderTimeout    time.Duration
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		HTTPPort:                 defaultHTTPPort,
		HTTPReadHeaderTimeout:    defaultHTTPReadHeaderTimeout,
		URLPathFirstSubdirectory: defaultURLPathFirstSubdirectory,
	}
}

// Profiler provides HTTP endpoints for pprof profiling.
type Profiler struct {
	config  *Config
	started bool
}

// New creates a new Profiler instance with the default configuration.
func New() *Profiler {
	return &Profiler{
		config: DefaultConfig(),
	}
}

// WithConfig sets a custom configuration for the profiler.
func (profiler *Profiler) WithConfig(config *Config) *Profiler {
	if config == nil {
		config = DefaultConfig()
	}
	profiler.config = config

	return profiler
}

// WithHTTPPort sets a custom HTTP port for the profiler.
func (profiler *Profiler) WithHTTPPort(port int) *Profiler {
	profiler.config.HTTPPort = port

	return profiler
}

// WithReadHeaderTimeout sets a custom read header timeout.
func (profiler *Profiler) WithReadHeaderTimeout(timeout time.Duration) *Profiler {
	profiler.config.HTTPReadHeaderTimeout = timeout

	return profiler
}

// WithHTTPServer sets a custom HTTP server.
func (profiler *Profiler) WithHTTPServer(server *http.Server) *Profiler {
	profiler.config.HTTPServer = server

	return profiler
}

// WithDefaultURLPathFirstSubdirectory sets a custom URL path prefix.
func (profiler *Profiler) WithDefaultURLPathFirstSubdirectory(firstSubdirectory string) *Profiler {
	profiler.config.URLPathFirstSubdirectory = firstSubdirectory

	return profiler
}

// Start initializes and starts the profiler server.
func (profiler *Profiler) Start(ctx context.Context) error {
	if profiler.started {
		return ErrServerAlreadyStarted
	}

	if profiler.config.HTTPServer == nil {
		profiler.config.HTTPServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", profiler.config.HTTPPort),
			ReadHeaderTimeout: profiler.config.HTTPReadHeaderTimeout,
		}
	}

	mux := http.NewServeMux()
	prefix := "/" + profiler.config.URLPathFirstSubdirectory

	// Register pprof handlers.
	mux.HandleFunc(prefix+"/profile", pprof.Profile)
	mux.HandleFunc(prefix+"/index", pprof.Index)
	mux.HandleFunc(prefix+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"/symbol", pprof.Symbol)
	mux.Handle(prefix+"/goroutine", pprof.Handler("goroutine"))
	mux.Handle(prefix+"/heap", pprof.Handler("heap"))
	mux.Handle(prefix+"/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle(prefix+"/block", pprof.Handler("block"))
	mux.Handle(prefix+"/trace", pprof.Handler("trace"))
	mux.Handle(prefix+"/allocs", pprof.Handler("allocs"))

	profiler.config.HTTPServer.Handler = mux
	profiler.started = true

	go func() {
		slog.Info("Starting profiler server", "port", profiler.config.HTTPPort)
		if err := profiler.config.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Profiler server error", "error", err)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down profiler server")
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := profiler.config.HTTPServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down profiler server", "error", err)
		}
		slog.Info("Profiler server shut down")
	}()

	return nil
}
