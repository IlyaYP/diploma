package server

import (
	"context"
	"fmt"
	"github.com/IlyaYP/diploma/api/server/handler"
	"github.com/IlyaYP/diploma/pkg/logging"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

const (
	serviceName = "http-server"
)

// Config provides the configuration for the API server
type Config struct {
	Address string
}

// Server contains instance details for the server
type (
	Server struct {
		*http.Server
		cfg *Config
	}
	Option func(s *Server) error
)

// New returns a new instance of the server based on the specified configuration.
// It allocates resources which will be needed for ServeAPI(ports, unix-sockets).
func New(opts ...Option) (*Server, error) {
	s := &Server{}
	s.Server = &http.Server{}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("initialising dependencies: %w", err)
		}
	}

	if s.cfg == nil {
		return nil, fmt.Errorf("config: nil")
	}

	s.Addr = s.cfg.Address

	return s, nil
}

// WithConfig sets Config.
func WithConfig(cfg *Config) Option {
	return func(s *Server) error {
		s.cfg = cfg
		return nil
	}
}

// WithRouter sets Router.
func WithRouter(r *handler.Handler) Option {
	return func(s *Server) error {
		s.Handler = r
		return nil
	}
}

// Serve starts listening for inbound requests.
func (s *Server) Serve(ctx context.Context) {
	ctx, _ = logging.GetCtxLogger(ctx) // correlationID is created here
	logger := s.Logger(ctx)
	logger.Info().Msg("Started serve connections")

	// service connections
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Err(err).Msg("ListenAndServe")
		return
	}
	logger.Info().Msg("Finished serve connections")
}

// Close closes the HTTPServer from listening for the inbound requests.
func (s *Server) Close(ctx context.Context) {
	logger := s.Logger(ctx)
	logger.Info().Msg("Shutdown Server")
	ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctxt); err != nil {
		logger.Err(err).Msg("HTTP server Shutdown")
	}
}

func (s *Server) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
