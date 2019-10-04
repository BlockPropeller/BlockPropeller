package server

import (
	"fmt"
	"net/http"
	"time"

	"chainup.dev/lib/log"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Server is a wrapper around the echo.Echo HTTP server that allows us to configure
// all the aspects of the server upfront, leaving us only with a Start() method to call.
type Server struct {
	httpSrv *http.Server
	echoSrv *echo.Echo
}

// ProvideServer configures a Server instance and prepares it for listening for new requests.
func ProvideServer(cfg *Config, router Router, logger log.Logger) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.Use(LoggerMiddleware(logger))

	err := router.RegisterRoutes(e)
	if err != nil {
		return nil, errors.Wrap(err, "register routes")
	}

	return &Server{
		httpSrv: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			ReadTimeout:  cfg.ReadTimeout * time.Second,
			WriteTimeout: cfg.WriteTimeout * time.Second,
		},
		echoSrv: e,
	}, nil
}

// Start the server and wait until an error is returned.
func (srv *Server) Start() error {
	log.Info("Started listening for HTTP requests", log.Fields{
		"addr": srv.httpSrv.Addr,
	})
	err := srv.echoSrv.StartServer(srv.httpSrv)
	if err != nil {
		return errors.Wrap(err, "start server")
	}
	return nil
}
