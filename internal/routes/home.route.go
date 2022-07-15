package routes

import (
	"github.com/AliTr404/T-MO/internal/http/controllers"
	"github.com/labstack/echo/v4"
)

func HomeRoutes(server *echo.Echo, homeController controllers.HomeController, middleware ...echo.MiddlewareFunc) {
	homeRoutes := server.Group("/api/home", middleware...)
	homeRoutes.GET("/search", homeController.Search)
	homeRoutes.GET("/latests", homeController.GetLatestsVideo)
}
