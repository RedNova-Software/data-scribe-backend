package util

import "time"

func GetCurrentTime() int64 {
	now := time.Now()
	return now.Unix()
}
