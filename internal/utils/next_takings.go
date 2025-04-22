package utils

import (
	"errors"
	"time"
)

func GenerateScheduleTimes(day time.Time, timesPerDay uint64) ([]time.Time, error) {
	if timesPerDay <= 0 || timesPerDay > 24 {
		return nil, errors.New("Неверное количество приемов в день")
	}

	var scheduleTimes []time.Time
	startTime := time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())
	now := time.Now()

	if now.After(endTime) {
		day = day.AddDate(0, 0, 1)
		startTime = time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
		endTime = time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())
	}

	if timesPerDay == 24 { // каждый час
		for i := 0; i < 24; i++ {
			appointmentTime := time.Date(day.Year(), day.Month(), day.Day(), i, 0, 0, 0, day.Location())
			scheduleTimes = append(scheduleTimes, appointmentTime)
		}
	} else {
		startTime = time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
		endTime = time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())

		interval := (endTime.Sub(startTime)) / time.Duration(timesPerDay)

		for i := 0; uint64(i) < timesPerDay; i++ {
			appointmentTime := startTime.Add(time.Duration(i) * interval)
			scheduleTimes = append(scheduleTimes, appointmentTime)
		}
	}

	return scheduleTimes, nil
}
