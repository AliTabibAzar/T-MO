package routes

import (
	"github.com/AliTr404/T-MO/internal/http/controllers"
	mw "github.com/AliTr404/T-MO/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

func UploadRoutes(server *echo.Echo, uploadController controllers.UploadController, middleware ...echo.MiddlewareFunc) {
	uploadRoutes := server.Group("/api/upload", middleware...)
	uploadRoutes.POST("/video", uploadController.UploadVideo, mw.UploadFile("./upload/private/", mw.MAX_VIDEO_FILE_SIZE, "", mw.VideoFileFilter))
	uploadRoutes.POST("/profile", uploadController.UploadProfilePicture, mw.UploadFile("./upload/public/", mw.MAX_IMAGE_FILE_SIZE, "/files", mw.ImageFileFilter))
}
