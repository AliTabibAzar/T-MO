package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AliTr404/T-MO/internal/dto"
	"github.com/AliTr404/T-MO/internal/http/exception"
	"github.com/AliTr404/T-MO/internal/services"
	"github.com/AliTr404/T-MO/pkg/derrors"
	"github.com/AliTr404/T-MO/pkg/tol"
	"github.com/AliTr404/T-MO/pkg/validation"
	"github.com/labstack/echo/v4"
)

type (
	VideoController interface {
		PlayVideo(c echo.Context) error
		UpdateVideo(c echo.Context) error
		LikeVideo(c echo.Context) error
		UnLikeVideo(c echo.Context) error
		SendComment(c echo.Context) error
		GetComments(c echo.Context) error
		GetCommentsReply(c echo.Context) error
	}

	videoController struct {
		videoService services.VideoService
	}
)

func NewVideoController(videoService services.VideoService) VideoController {
	return &videoController{
		videoService: videoService,
	}
}

func (s *videoController) PlayVideo(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => PlayVideo: %v", c.RealIP()))

	video, err := s.videoService.PlayVideo(c.Param("id"))
	if err != nil {
		return derrors.DHttpError(err)
	}
	file, err := os.Open(video.Path)
	if err != nil {
		return exception.InternalServerException()
	}
	return c.Stream(http.StatusOK, "video/mp4", file)
}
func (s *videoController) UpdateVideo(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => UpdateVideo: %v", c.RealIP()))

	data := new(dto.Video)
	if err := c.Bind(&data); err != nil {
		return exception.BadRequestException("اطلاعات وارد شده دارای خطا است")
	}
	if err := validation.ValidateRequest(c, data); err != nil {
		return exception.BadRequestException(err)
	}
	videoID := c.Param("id")
	if err := s.videoService.UpdateVideo(videoID, *data); err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ویدیو با موفقیت بروزرسانی شد.",
	})
}
func (s *videoController) LikeVideo(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => LikeVideo: %v", c.RealIP()))

	userID := c.Get("userID").(string)
	videoID := c.Param("id")
	if err := s.videoService.LikeVideo(videoID, userID); err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ویدیو با موفقیت لایک شد.",
	})
}
func (s *videoController) UnLikeVideo(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => UnLikeVideo: %v", c.RealIP()))

	userID := c.Get("userID").(string)
	videoID := c.Param("id")
	if err := s.videoService.UnLikeVideo(videoID, userID); err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ویدیو با موفقیت آنلایک شد.",
	})
}

func (s *videoController) SendComment(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => SendComment: %v", c.RealIP()))

	data := new(dto.Comment)
	if err := c.Bind(&data); err != nil {
		return exception.BadRequestException("اطلاعات وارد شده دارای خطا است")
	}
	if err := validation.ValidateRequest(c, data); err != nil {
		return exception.BadRequestException(err)
	}
	userID := c.Get("userID").(string)
	videoID := c.Param("id")
	err := s.videoService.SendComment(videoID, userID, data)
	if err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "کامنت با موفقیت ثبت شد.",
	})
}

func (s *videoController) GetComments(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => GetComments: %v", c.RealIP()))

	videoID := c.Param("id")
	comments, err := s.videoService.GetComments(videoID)
	if err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": len(comments),
		"data":  comments,
	})
}

func (s *videoController) GetCommentsReply(c echo.Context) error {
	tol.TMessage(fmt.Sprintf("Controller (Video) => GetCommentsReply: %v", c.RealIP()))

	commentID := c.Param("comment_id")
	comments, err := s.videoService.GetCommentsReply(commentID)
	if err != nil {
		return derrors.DHttpError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": len(comments),
		"data":  comments,
	})
}
