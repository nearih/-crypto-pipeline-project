package http

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo"
)

type HttpServer struct {
	Echo *echo.Echo
}

func NewHttpServer() *HttpServer {
	e := echo.New()

	return &HttpServer{e}
}

func (sv *HttpServer) Start(port int) {
	if port == 0 {
		port = 9090
	}
	sv.Echo.Start(fmt.Sprintf(":%d", port))
}

func (sv *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sv.Echo.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
