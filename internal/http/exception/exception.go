package exception

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.JSON(report.Code, map[string]interface{}{
		"message": report.Message,
		"Code":    report.Code,
	})
}

func ConflictException(msg interface{}) error {
	return echo.NewHTTPError(http.StatusConflict, msg)
}
func BadRequestException(msg interface{}) error {
	return echo.NewHTTPError(http.StatusBadRequest, msg)
}
func UnauthorizedException(msg interface{}) error {
	return echo.NewHTTPError(http.StatusUnauthorized, msg)
}
func NotFoundException(msg interface{}) error {
	return echo.NewHTTPError(http.StatusNotFound, msg)
}
func MethodNotAllowedException(msg interface{}) error {
	return echo.NewHTTPError(http.StatusMethodNotAllowed, msg)
}
func RequestEntityTooLargeException(msg interface{}) error {
	return echo.NewHTTPError(http.StatusRequestEntityTooLarge, msg)
}
func InternalServerException() error {
	return echo.NewHTTPError(http.StatusInternalServerError, "خطای سرور")
}
