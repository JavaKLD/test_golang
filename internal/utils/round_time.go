package utils

import "time"

func RoundTime(t time.Time) time.Time {
	minRound := ((t.Minute() + 14) / 15) * 15
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minRound%60, 0, 0, t.Location())
}
