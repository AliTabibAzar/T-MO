package thumbnail

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/AliTr404/T-MO/pkg/SID"
)

type Thumbnail interface {
	Generate() (string, error)
}
type thumbnail struct {
	filepath string
	width    uint16
	height   uint16
}

func NewThumbnail(filepath string, width uint16, height uint16) Thumbnail {
	return &thumbnail{
		filepath: filepath,
		width:    width,
		height:   height,
	}
}
func (t *thumbnail) Generate() (string, error) {
	filename := SID.SIDgenerator(20) + ".png"
	output := "./upload/public/" + filename
	cmd := exec.Command("ffmpeg", "-i", t.filepath, "-ss", "00:00:00.00", "-vframes", "1", "-s", fmt.Sprintf("%dx%d", t.width, t.height), output)
	if cmd.Run() != nil {
		fmt.Print(cmd.Run().Error())
		return "", errors.New("failed to generate thumbnail")
	}
	return "/files/" + filename, nil
}
