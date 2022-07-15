package controllers

import (
	"net/http"

	"github.com/AliTr404/T-MO/internal/services"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/labstack/echo/v4"
)

type (
	HomeController interface {
		Search(c echo.Context) error
		GetLatestsVideo(c echo.Context) error
	}

	homeController struct {
		homeService services.HomeService
	}
)

func NewHomeController(homeService services.HomeService) HomeController {
	return &homeController{
		homeService: homeService,
	}
}

func (s *homeController) Search(c echo.Context) error {
	results, err := s.homeService.Search(c.QueryParams())
	if err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": len(results),
		"data":  results,
	})
}

func (s *homeController) GetLatestsVideo(c echo.Context) error {
	results, err := s.homeService.GetLatestsVideo(c.QueryParams())
	if err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": len(results),
		"data":  results,
	})
}
