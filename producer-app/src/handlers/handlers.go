package handlers

import (
	"context"
	"fmt"
	"producer-app/server"
	"producer-app/src/services"

	"github.com/labstack/echo"
)

type Handlers struct {
	Ctx        context.Context
	Services   *services.Services
	HttpServer *server.HttpServer
}

func NewHandlers(ctx context.Context, e *server.HttpServer, s *services.Services) *Handlers {

	h := &Handlers{
		Ctx:        ctx,
		Services:   s,
		HttpServer: e,
	}

	h.registerHandler()

	return h
}

func (h *Handlers) registerHandler() {
	h.HttpServer.Echo.GET("/ping", h.Ping())
	h.HttpServer.Echo.GET("/getdata", h.GetData())
}

func (h *Handlers) Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, "success")
	}
}

func (h *Handlers) GetData() echo.HandlerFunc {
	return func(c echo.Context) error {
		symbol := c.QueryParam("symbol")
		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		go func() {
			<-h.Ctx.Done()
			cancel()
		}()

		err := h.Services.UploadData(ctx, symbol)
		if err != nil {
			return c.JSON(500, fmt.Sprintf("error: %v", err))
		}
		return c.JSON(200, fmt.Sprintf("get data from %v", symbol))
	}

}
