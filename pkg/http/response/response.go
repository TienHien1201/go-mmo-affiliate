// pkg/http/response/response.go
package response

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

// dataResponse là internal builder — không export
func dataResponse(c echo.Context, statusCode int, data any, opts []ResponseOption) error {
    o := defaultOptions()
    o.apply(opts)

    return c.JSON(statusCode, SuccessResDto{
        Status:  statusCode,
        Message: http.StatusText(statusCode),
        Data:    data,
        Doc:     o.doc,
        Meta:    o.meta,
    })
}

// errResponse là internal builder cho error — không export
func errResponse(c echo.Context, statusCode int, data any, opts []ResponseOption) error {
    o := defaultOptions()
    o.apply(opts)

    return c.JSON(statusCode, ErrorResDto{
        Status:  statusCode,
        Message: http.StatusText(statusCode),
        Code:    o.code,
        Data:    data,
        Doc:     o.doc,
        Meta:    o.meta,
    })
}