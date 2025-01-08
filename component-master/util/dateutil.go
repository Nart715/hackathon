package util

import "time"

func MilliToTime(t int64) *time.Time {
	resp := time.UnixMilli(t)
	return &resp
}
