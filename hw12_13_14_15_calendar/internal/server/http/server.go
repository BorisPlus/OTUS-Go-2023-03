package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	api "hw12_13_14_15_calendar/internal/server/http/api"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	server *http.Server
	// host   string
	// port   uint16
	logger interfaces.Logger
	app    interfaces.Applicationer
	cancel context.CancelFunc
}

func NewServer(
	host string,
	port uint16,
	readTimeout time.Duration, // TODO: set default in ServerConfig `10 * time.Second`
	readHeaderTimeout time.Duration, // TODO: set default in ServerConfig `10 * time.Second`
	writeTimeout time.Duration, // TODO: set default in ServerConfig `10 * time.Second`
	maxHeaderBytes int, // TODO: set default in ServerConfig `1 << 20`
	logger interfaces.Logger,
	app interfaces.Applicationer,
) *Server {
	server := http.Server{
		Addr:              net.JoinHostPort(host, fmt.Sprint(port)),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	return &Server{&server, logger, app, nil}
}

func (s *Server) Start(ctx context.Context) error {
	// _ = ctx // TODO: for what?
	s.logger.Info("Server.Start()")

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Init(s.logger).Listen(http.HandlerFunc(handleTeapot)))
	mux.Handle("/api",  api.Handlers(s.logger, s.app))
	s.server.Handler = mux

	// // http.Handle("/", instance.Listen(http.HandlerFunc(handleTeapot)))
	// // http.Handle("/api/", api.Handlers(s.logger, s.app))
	// err := http.ListenAndServe(s.Address(), nil)
	// if err != nil {
	// 	return err
	// }
	// return nil

	contextHTTP, cancelHTTP := context.WithCancel(ctx)
	s.cancel = cancelHTTP
	s.server.BaseContext = func(l net.Listener) context.Context {return contextHTTP}

	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil

}

func handleTeapot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("I receive teapot-status code!"))
}

func (s *Server) Stop(ctx context.Context) error {
	// _ = ctx // TODO: for what?
	s.logger.Info("Server.Stop()")
	s.server.Shutdown(ctx)
	return nil
}
