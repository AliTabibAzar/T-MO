package derrors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	kind uint

	serverError struct {
		Kind    kind
		message string
	}
)

const (
	_ kind = iota
	KindInvalid
	KindNotFound
	KindConflict
	KindUnauthorized
	KindUnexpected
	KindNotAllowd
)

var (
	httpErrors = map[kind]int{
		KindInvalid:      http.StatusBadRequest,
		KindNotFound:     http.StatusNotFound,
		KindConflict:     http.StatusConflict,
		KindUnauthorized: http.StatusUnauthorized,
		KindUnexpected:   http.StatusInternalServerError,
		KindNotAllowd:    http.StatusMethodNotAllowed,
	}
)

func New(kind kind, message string) error {
	return &serverError{
		Kind:    kind,
		message: message,
	}
}

func (s *serverError) Error() string {
	return s.message
}

func DHttpError(err error) *echo.HTTPError {
	s, ok := err.(*serverError)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, s.message)
	}
	c, ok := httpErrors[s.Kind]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, s.message)
	}
	return echo.NewHTTPError(c, s.message)
}
