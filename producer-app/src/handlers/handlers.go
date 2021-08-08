package handlers

import (
	"context"
	"fmt"
	"net/http"
	"producer-app/server"
	"producer-app/src/model"
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
	h.HttpServer.Echo.POST("/getdata", h.GetData())
}

func (h *Handlers) Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, "success")
	}
}

func (h *Handlers) GetData() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.Request
		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		go func() {
			<-h.Ctx.Done()
			cancel()
		}()

		err = h.Services.UploadData(ctx, req.Symbol)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
		}
		return c.JSON(http.StatusOK, fmt.Sprintf("get data from %v", req.Symbol))
	}

}
