package tol

import (
	"fmt"
	"io"
	"os"

	"github.com/AliTr404/T-MO/pkg/dtime"
)

type Color uint

const (
	_ Color = iota
	Green
	Yellow
	Red
	Cyan
	White
	Blue
	Reset
)

var writer io.Writer = os.Stdout

var (
	colorValue = map[Color]string{
		Green:  "\033[92m",
		Yellow: "\033[93m",
		Red:    "\033[91m",
		Cyan:   "\033[92m",
		White:  "\033[97m",
		Blue:   "\033[94m",
		Reset:  "\033[0m",
	}
)

func (c Color) String() string {
	if s, ok := colorValue[c]; ok {
		return s
	}
	return fmt.Sprintf("ColorValue(#{b})")
}

func Init(w io.Writer) error {
	if w != nil {

		mw := io.MultiWriter(os.Stdout, w)
		writer = mw
		return nil
	}
	return fmt.Errorf("file argument has some issus")

}

func TWarning(s ...interface{}) {
	fmt.Fprintf(writer, "%v [TMO]-(Warning) | %v %v %v \n", Yellow.String(), dtime.TimeFormat(), s, Reset.String())
}
func TError(s ...interface{}) {
	fmt.Fprintf(writer, "%v [TMO]-(Error) | %v %v %v \n", Red.String(), dtime.TimeFormat(), s, Reset.String())
}
func TInfo(s ...interface{}) {
	fmt.Fprintf(writer, "%v [TMO]-(Info) | %v %v %v \n", Blue.String(), dtime.TimeFormat(), s, Reset.String())
}
func TMessage(s ...interface{}) {
	fmt.Fprintf(writer, "%v [TMO]-(Message) | %v %v %v \n", Cyan.String(), dtime.TimeFormat(), s, Reset.String())
}
