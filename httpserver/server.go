package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"
)

// Server is a wrapper around the http.Server, storing
// the logger, and options
type Server struct {
	options *Options
	logger  *zap.Logger
	server  *http.Server
}

// New instantiates a new httpserver with the options provided
func New(logger *zap.Logger, options *Options) *Server {
	return &Server{
		options: options,
		logger:  logger.With(zap.String("module", "httpserver")),
	}
}

// Run runs the http server with the http.DefaultServeMux
// in a seperate goroutine until Shutdown is run
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.options.listenAddr)
	if err != nil {
		return fmt.Errorf("could not run http server: %w", err)
	}
	s.logger.
		With(zap.String("address", s.options.listenAddr)).
		Info("listening")

	s.server = &http.Server{
		Handler: http.DefaultServeMux,
	}
	go func() {
		err := s.server.Serve(listener)
		if err != http.ErrServerClosed {
			s.logger.
				With(zap.Error(err)).
				Error("httpserver encountered an error")
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the http server -
// immediately stops accepting new connections, but
// continues processing existing ones, then terminating
// when all existing connections are closed, or there
// is a timeout
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.
		With(zap.Duration("duration", s.options.shutdownDuration)).
		Info("shutdown")

	ctx, cancel := context.WithTimeout(ctx, s.options.shutdownDuration)
	defer cancel()

	return s.server.Shutdown(ctx)
}
