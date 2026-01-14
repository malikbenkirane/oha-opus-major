package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/malikbenkirane/oha-opus-major/internal/port"
)

// Config holds server configuration parameters.
type Config struct {
	// addr is the address on which the server will listen, e.g., ":8080".
	addr string
	// readTimeout is the maximum duration for reading the entire request, including the body.
	readTimeout time.Duration
	// writeTimeout is the maximum duration before timing out writes of the response.
	writeTimeout time.Duration
	// idleTimeout is the maximum amount of time to wait for the next request when keep‑alives are enabled.
	idleTimeout time.Duration
	// shutdownTimeout defines the maximum duration to wait for the server to shut down gracefully.
	shutdownTimeout time.Duration
}

// newServer creates and returns a pointer to an http.Server using the values from Config.
func (c Config) newServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         c.addr,
		ReadTimeout:  c.readTimeout,
		WriteTimeout: c.writeTimeout,
		IdleTimeout:  c.idleTimeout,
		Handler:      mux,
	}
}

// defaultConfig provides example defaults for illustration purposes.
// Update these settings to suit a production multiplayer player‑data server.
func defaultConfig(addr string) Config {
	return Config{
		addr:            addr,
		readTimeout:     500 * time.Millisecond,
		writeTimeout:    1 * time.Second,
		idleTimeout:     15 * time.Second,
		shutdownTimeout: 15 * time.Second,
	}
}

func New(addr string, opts ...Option) port.PlayerDataServer {
	conf := defaultConfig(addr)
	for _, opt := range opts {
		conf = opt(conf)
	}

	server := &server{
		config: conf,
		err:    make(chan error, 1),
	}

	mux := http.NewServeMux()
	mux.Handle("GET /update-player-data", server.handler(server.handleGetUpdatePlayerData))

	server.server = conf.newServer(mux)

	return server
}

type server struct {
	server *http.Server
	config Config
	err    chan error
}

// Serve starts the HTTP server in a goroutine, monitors the provided
// context for cancellation, and performs a graceful shutdown using the
// configured timeout. Any error returned by ListenAndServe is wrapped with
// the listening address and propagated to the caller.
func (s server) Serve(ctx context.Context) error {
	// Monitor server errors
	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Debug("server error handler loop: done with context")
				return
			case err := <-s.err:
				// Simple error handling for this assignment.
				slog.Error("server error", "err", err)
				return
			}
		}
	}()

	// Capture the result of ListenAndServe.
	srvErr := make(chan error, 1)
	defer close(srvErr)
	go func() {
		srvErr <- s.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		// Context cancelled – initiate graceful shutdown.
		shutdownCtx, cancel := context.WithTimeout(ctx, s.config.shutdownTimeout)
		defer cancel()
		return s.shutdown(shutdownCtx)

	case err, ok := <-srvErr:
		if !ok {
			return fmt.Errorf("reading on closed channel srvErr")
		}
		return fmt.Errorf("listen and serve on %q: %w", s.config.addr, err)
	}
}

func (s server) shutdown(ctx context.Context) error {
	close(s.err)
	return s.server.Shutdown(ctx)
}

type handler func(w http.ResponseWriter, r *http.Request) error

func (s server) handler(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// For this assignment we keep the error handling straightforward.
			s.err <- fmt.Errorf("handle %s %q: %w", r.Method, r.URL, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				s.err <- fmt.Errorf("handle %s %q: %w", r.Method, r.URL, err)
			}
		}
	}
}

var ErrNotImplemented = errors.New("not implemented")

func (s server) handleGetUpdatePlayerData(w http.ResponseWriter, r *http.Request) error {
	// TODO
	return ErrNotImplemented
}

type Option func(Config) Config

// WithReadTimeout configures the HTTP server's read timeout.
// It returns an Option that updates the Config's readTimeout field.
func WithReadTimeout(timeout time.Duration) Option {
	return func(c Config) Config {
		c.readTimeout = timeout
		return c
	}
}

// WithWriteTimeout configures the HTTP server's write timeout.
// It returns an Option that updates the Config's writeTimeout field.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(c Config) Config {
		c.writeTimeout = timeout
		return c
	}
}

// WithIdleTimeout configures the HTTP server's idle timeout.
// It returns an Option that updates the Config's idleTimeout field.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(c Config) Config {
		c.idleTimeout = timeout
		return c
	}
}

// WithShutdow configures the HTTP server's graceful shutdown timeout.
// It returns an Option that updates the Config's shutdownTimeout field.
func WithShutdowTimeout(timeout time.Duration) Option {
	return func(c Config) Config {
		c.shutdownTimeout = timeout
		return c
	}
}
