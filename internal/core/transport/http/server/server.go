package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Slazzzer/golang-todoapp/docs"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) Mux() *http.ServeMux {
	return s.mux
}

func (s *HTTPServer) Config() Config {
	return s.config
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		handler := core_http_middleware.ChainMiddleware(router, router.middleware...)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, handler))
	}
}

func (s *HTTPServer) RegisterSwaggerRouter() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:              s.config.Addr,
		Handler:           mux,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("Starting HTTP server", core_logger.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP server: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("Shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server shutdown")
	}
	return nil
}
