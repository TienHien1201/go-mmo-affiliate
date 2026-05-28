package handler

import (
	xhttp "github.com/TienHien1201/go-mmo-affiliate/pkg/http"
	xres "github.com/TienHien1201/go-mmo-affiliate/pkg/http/response"
	xlogger "github.com/TienHien1201/go-mmo-affiliate/pkg/logger"
	echoSwagger "github.com/swaggo/echo-swagger"

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

// HealthCheck godoc
// @Summary Health check
// @Description Check API status
// @Tags System
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /api/v1/health [get]
func (h *handler) HealthCheck(c echo.Context) error {
	return xres.Success(c, nil,xres.WithDoc("/api/v1/docs/index.html"))
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	api.GET("/health", h.HealthCheck)

	// swagger endpoint
	api.GET("/docs/*", echoSwagger.WrapHandler)
}
