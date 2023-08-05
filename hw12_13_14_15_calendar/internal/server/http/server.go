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
	server      *http.Server
	logger      interfaces.Logger
	app         interfaces.Applicationer
	cancel      context.CancelFunc
	// baseContext context.Context
	context     context.Context
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

	logger.Info(net.JoinHostPort(host, fmt.Sprint(port)))
	server := http.Server{
		Addr:              net.JoinHostPort(host, fmt.Sprint(port)),
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	return &HTTPServer{&server, logger, app, nil, nil}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	errChannel := make(chan error)
	s.logger.Info("Server.Start()")
	s.context, s.cancel = context.WithCancel(ctx)
	select {
	case <-(ctx).Done(): // Почему не работает при cancel в родителе
		fmt.Println("<-(*ctx).Done()")
		err := s.Stop()
		if err != nil {
			return err
		}
		return nil
	case <-s.context.Done():
		fmt.Println("<-s.context.Done()")
		err := s.Stop()
		if err != nil {
			return err
		}
		return nil
	case errChannel <- s.server.ListenAndServe():
		fmt.Println("s.server.ListenAndServe()")
		err := <-errChannel
		s.logger.Error(err.Error())
		return err
	}
}

func (s *HTTPServer) Stop() error {
	s.logger.Info("Server.Stop()")
	if s.context != nil {
		s.cancel()
		err := s.server.Shutdown(s.context)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleTeapot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("I receive teapot-status code!"))
}
