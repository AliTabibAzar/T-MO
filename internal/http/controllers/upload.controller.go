package controllers

import (
	"net/http"

	"github.com/AliTr404/T-MO/internal/models"
	"github.com/AliTr404/T-MO/internal/services"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UploadController interface {
		UploadVideo(c echo.Context) error
		UploadProfilePicture(c echo.Context) error
	}

	uploadController struct {
		uploadService services.UploadService
	}
)

func NewUploadController(uploadService services.UploadService) UploadController {
	return &uploadController{
		uploadService: uploadService,
	}
}

func (s *uploadController) UploadVideo(c echo.Context) error {
	caption := c.FormValue("caption")
	path := c.Get("filePath").(string)
	userid, _ := primitive.ObjectIDFromHex(c.Get("userID").(string))
	video, err := s.uploadService.UploadVideo(&models.Video{Path: path, Caption: caption, User: userid})
	if err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ویدیو با موفقیت آپلود شد.",
		"data":    video,
	})
}

func (s *uploadController) UploadProfilePicture(c echo.Context) error {
	path := c.Get("filePath").(string)

	if err := s.uploadService.UploadProfilePicture(c.Get("userID").(string), path); err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "عکس پروفایل با موفقیت آپلود شد.",
	})
}
