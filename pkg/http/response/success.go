package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, data any, opts ...ResponseOption) error {
    return dataResponse(c, http.StatusOK, data, opts)
}

func Created(c echo.Context, data any, opts ...ResponseOption) error {
    return dataResponse(c, http.StatusCreated, data, opts)
}

func NoContent(c echo.Context) error {
    return c.NoContent(http.StatusNoContent)
}