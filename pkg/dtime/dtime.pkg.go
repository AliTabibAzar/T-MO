package dtime

import "time"

func Now() time.Time {
	return time.Now()
}

func TimeFormat() string {
	return Now().Format("2006/01/02 15:04:05")
}
