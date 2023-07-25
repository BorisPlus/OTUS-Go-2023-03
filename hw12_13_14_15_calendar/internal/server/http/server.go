package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
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
	http.Handle("/", middleware(http.HandlerFunc(handleTeapot), s.logger))
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
