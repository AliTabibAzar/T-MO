package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AliTr404/T-MO/internal/dto"
	"github.com/AliTr404/T-MO/internal/http/exception"
	"github.com/AliTr404/T-MO/internal/models"
	"github.com/AliTr404/T-MO/internal/services"
	"github.com/AliTr404/T-MO/pkg/cookie"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/AliTr404/T-MO/pkg/jwt"
	"github.com/AliTr404/T-MO/pkg/tol"
	"github.com/AliTr404/T-MO/pkg/validation"
	"github.com/labstack/echo/v4"
)

type (
	AuthController interface {
		SignUp(c echo.Context) error
		SignIn(c echo.Context) error
		ForgotPassword(c echo.Context) error
		ChangePassword(c echo.Context) error
	}

	authController struct {
		authService services.AuthService
	}
)

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (s *authController) SignUp(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Auth) => SignUp: %v", c.RealIP()))

	data := new(models.User)
	if err := c.Bind(&data); err != nil {
		return exception.BadRequestException("اطلاعات وارد شده دارای خطا است")
	}

	if err := validation.ValidateRequest(c, data); err != nil {
		return exception.BadRequestException(err)
	}
	res, err := s.authService.SignUp(data)

	if err != nil {
		return derrors.DHttpError(err)
	}

	accessToken := jwt.GenerateAccessToken(res.ID.Hex())
	refreshToken := jwt.GenerateRefreshToken(res.ID.Hex())

	cookie.NewCookie("refreshToken", refreshToken, 730*time.Hour, true).SetCookie(c)
	cookie.NewCookie("accessToken", accessToken, 1*time.Hour, true).SetCookie(c)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": map[string]interface{}{
			"refresh_token": refreshToken,
			"access_token":  accessToken,
		},
		"user": res,
	})
}

func (s *authController) SignIn(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Auth) => SignIn: %v", c.RealIP()))

	data := new(dto.SignIn)
	if err := c.Bind(&data); err != nil {
		return exception.BadRequestException("اطلاعات وارد شده دارای خطا است")
	}

	if err := validation.ValidateRequest(c, data); err != nil {
		return exception.BadRequestException(err)
	}
	res, err := s.authService.SignIn(data)

	if err != nil {
		return derrors.DHttpError(err)
	}

	accessToken := jwt.GenerateAccessToken(res.ID.Hex())
	refreshToken := jwt.GenerateRefreshToken(res.ID.Hex())

	cookie.NewCookie("refreshToken", refreshToken, 730*time.Hour, true).SetCookie(c)
	cookie.NewCookie("accessToken", accessToken, 1*time.Hour, true).SetCookie(c)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": map[string]interface{}{
			"refresh_token": refreshToken,
			"access_token":  accessToken,
		},
		"user": res,
	})
}

func (s *authController) ForgotPassword(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Auth) => ForgotPassword: %v", c.RealIP()))

	data := new(dto.ForgotPassword)
	if err := c.Bind(&data); err != nil {
		return exception.BadRequestException("اطلاعات وارد شده دارای خطا است")
	}

	if err := validation.ValidateRequest(c, data); err != nil {
		return exception.BadRequestException(err)
	}
	if err := s.authService.ForgotPassword(data.Email); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "لینک بازیابی رمز عبور با موفقیت به ایمیل شما ارسال شد.",
	})
}

func (s *authController) ChangePassword(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Auth) => ChangePassword: %v", c.RealIP()))
	token := c.Param("token")
	if _, err := jwt.NewJwtManager(os.Getenv("JWT_FORGOT_PASSWORD_SECRET_KEY")).Verify(token); err != nil {
		return exception.UnauthorizedException([]dto.ApiError{{Param: "Token", Message: "توکن وارد شده معتبر نمی باشد !"}})
	}

	data := new(dto.ChangePassword)
	if err := c.Bind(&data); err != nil {
		return exception.BadRequestException("اطلاعات وارد شده دارای خطا است")
	}
	if err := validation.ValidateRequest(c, data); err != nil {
		return exception.BadRequestException(err)
	}
	if err := s.authService.ChangePassword(data); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "رمز عبور شما با موفقیت تغییر یافت.",
	})

}
