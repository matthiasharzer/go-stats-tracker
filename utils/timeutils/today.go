package timeutils

import "time"

func TodayDate() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}
