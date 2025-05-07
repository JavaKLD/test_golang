package utils

import "time"

func RoundTime(t time.Time) time.Time {
	minRound := ((t.Minute() + 14) / 15) * 15
	newHour := t.Hour()
	if minRound >= 60 {
		minRound = 0
		newHour++
	}
	return time.Date(t.Year(), t.Month(), t.Day(), newHour, minRound, 0, 0, t.Location())
}
