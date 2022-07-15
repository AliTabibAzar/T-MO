package routes

import (
	"github.com/AliTr404/T-MO/internal/http/controllers"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(server *echo.Echo, authController controllers.AuthController, middleware ...echo.MiddlewareFunc) {
	authRoutes := server.Group("/api/auth", middleware...)
	authRoutes.POST("/signup", authController.SignUp)
	authRoutes.POST("/signin", authController.SignIn)
}
