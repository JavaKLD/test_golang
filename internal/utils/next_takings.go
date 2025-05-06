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

	startTime = time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
	endTime = time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())

	var interval time.Duration
	if timesPerDay == 1 {
		interval = 0
	}
	if timesPerDay%2 == 0 {
		interval = (endTime.Sub(startTime)) / time.Duration(timesPerDay)
	}
	if timesPerDay%2 != 0 && timesPerDay != 1 {
		interval = (endTime.Sub(startTime)) / time.Duration(timesPerDay-1)
	}
	if timesPerDay == 24 {
		interval = endTime.Sub(startTime)
		for i := 0; uint64(i) <= 14; i++ {
			scheduleTimes = append(scheduleTimes, startTime.Add(time.Duration(i)*time.Hour))
		}
		return scheduleTimes, nil
	}

	for i := 0; uint64(i) < timesPerDay; i++ {
		appointmentTime := startTime.Add(time.Duration(i) * interval)
		scheduleTimes = append(scheduleTimes, RoundTime(appointmentTime))
	}

	return scheduleTimes, nil
}
