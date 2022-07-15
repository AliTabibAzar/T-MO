package routes

import (
	"github.com/AliTr404/T-MO/internal/http/controllers"
	"github.com/labstack/echo/v4"
)

func VideoRoutes(server *echo.Echo, videoController controllers.VideoController, middleware ...echo.MiddlewareFunc) {
	videoRoutes := server.Group("/api/video", middleware...)
	videoRoutes.GET("/:id", videoController.PlayVideo)
	videoRoutes.PUT("/:id", videoController.UpdateVideo)
	videoRoutes.GET("/:id/like", videoController.LikeVideo)
	videoRoutes.GET("/:id/unlike", videoController.UnLikeVideo)
	videoRoutes.GET("/:id/comments", videoController.GetComments)
	videoRoutes.POST("/:id/comments", videoController.SendComment)
	videoRoutes.GET("/:id/comments/:comment_id", videoController.GetCommentsReply)
}
