package xhttp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func DataResponse(c echo.Context, statusCode int, data interface{}, doc string) error {
	return c.JSON(statusCode, SuccessResDto{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
		Data:    data,
		Doc:     doc,
	})
}

func SuccessResponse(c echo.Context, data interface{}, doc string) error {
	return DataResponse(c, http.StatusOK, data, doc)
}

func CreatedResponse(c echo.Context, data interface{}, doc string) error {
	return DataResponse(c, http.StatusCreated, data, doc)
}

func NoContentResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func InternalServerErrorResponse(c echo.Context) error {
	return DataResponse(c, http.StatusInternalServerError, "Something went wrong", "")
}

// func BadRequestResponse(c echo.Context, data interface{}) error {
// 	return DataResponse(c, http.StatusBadRequest, data)
// }

// func UnauthorizedResponse(c echo.Context, data interface{}) error {
// 	return DataResponse(c, http.StatusUnauthorized, data)
// }
