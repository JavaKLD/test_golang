package utils

import "time"

const (
	roundingInterval = 15
	maxMinutesInHour = 60
)

func RoundTime(t time.Time) time.Time {
	minRound := ((t.Minute() + (roundingInterval - 1)) / roundingInterval) * roundingInterval
	newHour := t.Hour()

	if minRound >= maxMinutesInHour {
		minRound = 0
		newHour++
	}
	return time.Date(t.Year(), t.Month(), t.Day(), newHour, minRound, 0, 0, t.Location())
}
