package middleware

import (
	"io"
	"os"
	"path/filepath"

	"github.com/AliTr404/T-MO/internal/http/exception"
	"github.com/AliTr404/T-MO/pkg/SID"
	"github.com/labstack/echo/v4"
)

type UploadFilter func(string) error

const (
	MB                  = 1 << 20
	MAX_VIDEO_FILE_SIZE = 10 * MB
	MAX_IMAGE_FILE_SIZE = 1 * MB
)

func VideoFileFilter(ext string) error {
	switch ext {
	case ".mp4", ".gif", ".mov", ".m4v":
		return nil
	default:
		return exception.MethodNotAllowedException("فایل های مجاز برای آپلود عبارتند از : mp4, gif, mov, m4v")
	}
}

func ImageFileFilter(ext string) error {
	switch ext {
	case ".png", ".jpg", ".tiff":
		return nil
	default:
		return exception.MethodNotAllowedException("فایل های مجاز برای آپلود عبارتند از : png, jpg, tiff")
	}
}

func UploadFile(path string, maxFileSize int64, customPath string, filter UploadFilter) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			file, err := c.FormFile("file")
			if err != nil {
				return exception.BadRequestException("لطفا فایل را به درستی وارد کنید !")
			}
			if maxFileSize < file.Size {
				return exception.RequestEntityTooLargeException("حجم فایل وارد شده باید کمتر از 10 mb باشد !")
			}
			if err != nil {
				return exception.BadRequestException("عدم شناسایی فایل !")
			}
			if err = filter(filepath.Ext(file.Filename)); err != nil {
				return err
			}
			src, err := file.Open()
			if err != nil {
				return exception.BadRequestException("خطا در آپلود ویدیو !")
			}
			defer src.Close()
			newFileName := SID.SIDgenerator(16) + "-" + file.Filename
			path := path + newFileName
			dst, err := os.Create(path)
			if err != nil {
				return exception.BadRequestException("خطا در آپلود ویدیو !")
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return exception.BadRequestException("خطا در آپلود ویدیو !")
			}
			if customPath != "" {
				c.Set("filePath", customPath+"/"+newFileName)
			} else {
				c.Set("filePath", path)
			}

			return hf(c)
		}
	}
}
