package cookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AliTr404/T-MO/pkg/tol"
	"github.com/labstack/echo/v4"
)

type (
	CookiePkg interface {
		SetCookie(c echo.Context) *http.Cookie
		GetCookie(c echo.Context) (*http.Cookie, error)
	}

	Cookie struct {
		Key      string
		Value    string
		Expire   time.Duration
		HttpOnly bool
	}
)

func NewCookie(key string, value string, expire time.Duration, httpOnly bool) CookiePkg {
	return &Cookie{
		Key:      key,
		Value:    value,
		Expire:   expire,
		HttpOnly: httpOnly,
	}
}

func (ck *Cookie) SetCookie(c echo.Context) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = ck.Key
	cookie.Value = ck.Value
	cookie.Expires = time.Now().Add(ck.Expire)
	if ck.HttpOnly {
		cookie.HttpOnly = ck.HttpOnly
	}
	c.SetCookie(cookie)
	tol.TInfo(fmt.Sprintf("Cookie (pkg) | SetCookie => %v", cookie.Name))

	return cookie
}

func (ck *Cookie) GetCookie(c echo.Context) (*http.Cookie, error) {
	cookie, err := c.Cookie(ck.Key)
	if err != nil {
		return nil, err
	}
	tol.TInfo(fmt.Sprintf("Cookie (pkg) | GetCookie => %v", cookie.Name))

	return cookie, nil
}
