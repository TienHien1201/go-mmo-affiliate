package handler

import (
	xhttp "github.com/TienHien1201/go-mmo-affiliate/pkg/http"
	xlogger "github.com/TienHien1201/go-mmo-affiliate/pkg/logger"
	"github.com/labstack/echo/v4"
)

type handler struct {
	logger *xlogger.Logger
}

func NewHTTPHandler(logger *xlogger.Logger) xhttp.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) HeathyCheck(c echo.Context) error {
	return c.String(200, "OK")
}

func (h *handler) RegisterRoutes(e *echo.Echo){
	api := e.Group("api/v1")

	api.GET("heath", h.HeathyCheck)
}
