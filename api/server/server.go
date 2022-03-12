package server

import (
	"fmt"
	"github.com/IlyaYP/diploma/api/server/handler"
	"net/http"
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
func (s *Server) Serve() error {
	return s.ListenAndServe()
}

// TODO: smth
//// Close closes the HTTPServer from listening for the inbound requests.
//func (s *Server) Close() error {
//	//	return s.server.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	return s.Shutdown(ctx)
//
//}
