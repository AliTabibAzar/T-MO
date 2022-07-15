package middleware

import (
	"strings"

	"github.com/AliTr404/T-MO/internal/http/exception"
	"github.com/AliTr404/T-MO/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func Authorization(jwtManager jwt.JwtManager) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			auth := strings.Split(header, " ")
			if len(auth) <= 1 {
				return exception.UnauthorizedException("احراز هویت با خطا مواجه شد.")
			}
			token := auth[1]
			if _, err := jwtManager.Verify(token); err != nil {
				return exception.UnauthorizedException("احراز هویت با خطا مواجه شد.")
			}
			payload, _ := jwt.GetTokenPayload(token)
			c.Set("userID", payload.UserID)
			return hf(c)
		}
	}
}
