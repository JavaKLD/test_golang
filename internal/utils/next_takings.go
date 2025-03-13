package utils

import (
	"errors"
	"time"
)

func GenerateScheduleTimes(day time.Time, timesPerDay int) ([]time.Time, error) {
	if timesPerDay <= 0 {
		return nil, errors.New("Неверное количество приемов в день")
	}

	startTime := time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())

	interval := (endTime.Sub(startTime)) / time.Duration(timesPerDay)
	var scheduleTimes []time.Time

	for i := 0; i < timesPerDay; i++ {
		appointmentTime := startTime.Add(time.Duration(i) * interval)
		scheduleTimes = append(scheduleTimes, appointmentTime)
	}

	return scheduleTimes, nil
}
