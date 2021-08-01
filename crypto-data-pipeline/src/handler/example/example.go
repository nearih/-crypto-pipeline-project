package example

import "github.com/labstack/echo"

type ExampleHandler struct {
	Echo *echo.Echo
}

func NewExampleHandler(e *echo.Echo) *ExampleHandler {
	return &ExampleHandler{
		Echo: e,
	}
}

func (h *ExampleHandler) RegisterExampleHandler() {
	h.Echo.GET("/ping", h.Ping())
	h.Echo.GET("/example", h.GetExample())
}

func (h *ExampleHandler) Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, "success")
	}
}

func (h *ExampleHandler) GetExample() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, "success")
	}
}
