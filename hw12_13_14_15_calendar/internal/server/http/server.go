package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	api "hw12_13_14_15_calendar/internal/server/http/api"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
)

type Server struct {
	host   string
	port   uint16
	logger interfaces.Logger
	app    interfaces.Applicationer
	cancel context.CancelFunc
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func NewServer(host string, port uint16, logger interfaces.Logger, app interfaces.Applicationer) *Server {
	return &Server{host, port, logger, app, nil}
}

func (s *Server) Start(ctx context.Context) error {
	_ = ctx // TODO: for what?
	s.logger.Info("Server.Start()")
	instance := middleware.Init(s.logger)
	http.Handle("/", instance.Listen(http.HandlerFunc(handleTeapot)))
	// http.Handle("/api/", middleware.Middleware(http.HandleFunc("/api/", api.Routers.Go), s.logger))
	// as HandleFunc
	// http.HandleFunc("/api/", api.Routers(s.logger, s.app).ServeHTTP)
	// as Handle
	http.Handle("/api/", api.Handlers(s.logger, s.app))
	err := http.ListenAndServe(s.Address(), nil)
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
	_ = ctx // TODO: for what?
	s.logger.Info("Server.Stop()")
	return nil
}
