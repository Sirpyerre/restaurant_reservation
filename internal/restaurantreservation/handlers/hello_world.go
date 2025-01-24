package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HelloWorldHandler struct{}

func NewHelloWorldHandler() *HelloWorldHandler {
	return &HelloWorldHandler{}
}

func (h *HelloWorldHandler) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}
