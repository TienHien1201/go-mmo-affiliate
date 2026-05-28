package response

import "github.com/labstack/echo/v4"
import "net/http"

func BadRequest(c echo.Context, data any, opts ...ResponseOption) error {
    return errResponse(c, http.StatusBadRequest, data, opts)
}

func Unauthorized(c echo.Context, data any, opts ...ResponseOption) error {
    return errResponse(c, http.StatusUnauthorized, data, opts)
}

func Forbidden(c echo.Context, data any, opts ...ResponseOption) error {
    return errResponse(c, http.StatusForbidden, data, opts)
}

func NotFound(c echo.Context, data any, opts ...ResponseOption) error {
    return errResponse(c, http.StatusNotFound, data, opts)
}

func InternalServerError(c echo.Context) error {
    return errResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
}