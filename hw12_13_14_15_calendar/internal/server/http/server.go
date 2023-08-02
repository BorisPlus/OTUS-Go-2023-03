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

type HTTPServer struct {
	server *http.Server
	logger interfaces.Logger
	app    interfaces.Applicationer
	cancel context.CancelFunc
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
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	return &HTTPServer{&server, logger, app, nil}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	// _ = ctx // TODO: for what?
	s.logger.Info("Server.Start()")

	// // http.Handle("/", instance.Listen(http.HandlerFunc(handleTeapot)))
	// // http.Handle("/api/", api.Handlers(s.logger, s.app))
	// err := http.ListenAndServe(s.Address(), nil)
	// if err != nil {
	// 	return err
	// }
	// return nil

	contextHTTP, cancelHTTP := context.WithCancel(ctx)
	s.cancel = cancelHTTP
	s.server.BaseContext = func(l net.Listener) context.Context { return contextHTTP }

	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil

}

func goHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("go-go code!"))
}
func handleTeapot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("I receive teapot-status code!"))
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	// _ = ctx // TODO: for what?
	s.logger.Info("Server.Stop()")
	s.server.Shutdown(ctx)
	return nil
}
