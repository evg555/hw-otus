package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/gorilla/mux"
)

type Server struct {
	logger Logger
	app    Application
	srv    *http.Server
	cfg    *config.Config
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id string, title string) error
}

func NewServer(cfg config.Config, logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
		cfg:    &cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	r := mux.NewRouter()

	addr := net.JoinHostPort(s.cfg.App.Host, s.cfg.App.Port)
	handler := &Handler{
		router: r,
		app:    s.app,
		Logger: s.logger,
	}

	r.HandleFunc("/", handler.HomeHandler).Methods(http.MethodGet)
	r.Use(s.loggingMiddleware)

	s.srv = &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 1 * time.Second,
	}

	s.logger.Info(fmt.Sprintf("server starting at %s", addr))

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("server stopping...")
	return s.srv.Shutdown(ctx)
}
