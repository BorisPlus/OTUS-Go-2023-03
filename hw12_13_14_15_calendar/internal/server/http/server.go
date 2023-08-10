package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	api "hw12_13_14_15_calendar/internal/server/http/api"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
)

type HTTPServer struct {
	server  *http.Server
	logger  interfaces.Logger
	app     interfaces.Applicationer
	context context.Context
}

func NewHTTPServer(
	host string,
	port uint16,
	readTimeout time.Duration, // TODO: set time.Duration default "10s" in ServerConfig `10 * time.Second`
	readHeaderTimeout time.Duration, // TODO: set default "10s" in ServerConfig `10 * time.Second`
	writeTimeout time.Duration, // TODO: set default "10s" in ServerConfig `10 * time.Second`
	maxHeaderBytes int, // TODO: set default in ServerConfig `1 << 20`
	logger interfaces.Logger,
	app interfaces.Applicationer,
) *HTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/api/", api.Handlers(logger, app))
	mux.Handle("/", middleware.Instance().Listen(http.HandlerFunc(handleTeapot)))
	server := http.Server{
		Addr:              net.JoinHostPort(host, fmt.Sprint(port)),
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	this := &HTTPServer{}
	this.server = &server
	this.logger = logger
	this.app = app
	return this
}

func (s *HTTPServer) Start(ctx context.Context) error {
	s.logger.Info("HTTPServer.Start()")
	s.context = ctx
	go func() {
		<-ctx.Done()
		s.logger.Info("HTTPServer - Graceful Shutdown")
		if err := s.Stop(); err != nil {
			s.logger.Info("Stop error: %v\n", err)
		}
	}()
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.logger.Info("Start error: %v\n", err)
		return err
	}
	return nil
}

func (s *HTTPServer) Stop() error {
	err := s.server.Shutdown(s.context)
	if err != nil {
		return err
	}
	s.logger.Info("HTTPServer.Stop()")
	return nil
}

func handleTeapot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("I receive teapot-status code!"))
}
