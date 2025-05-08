package utils

import (
	"errors"
	"time"
)

var NowFunc = time.Now

func GenerateScheduleTimes(day time.Time, timesPerDay uint64) ([]time.Time, error) {
	if timesPerDay <= 0 || timesPerDay > 24 {
		return nil, errors.New("Неверное количество приемов в день")
	}

	var scheduleTimes []time.Time
	createdService := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())
	duration := endTime.Sub(startTime)
	now := NowFunc()

	if day.Year() < createdService.Year() || day.Year() > now.Year() {
		return nil, errors.New("неверный год")
	}

	if day.IsZero() {
		return nil, errors.New("некорректная дата: zero time")
	}

	if now.After(endTime) {
		day = day.AddDate(0, 0, 1)
		startTime = time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, day.Location())
		endTime = time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, day.Location())
	}

	if timesPerDay == 1 {
		scheduleTimes = append(scheduleTimes, startTime.Add(duration/2))
		return scheduleTimes, nil
	}

	if timesPerDay == 24 {
		for i := 0; uint64(i) <= 14; i++ {
			scheduleTimes = append(scheduleTimes, startTime.Add(time.Duration(i)*time.Hour))
		}
		return scheduleTimes, nil
	}

	interval := duration / time.Duration(timesPerDay-1)

	for i := 0; uint64(i) < timesPerDay; i++ {
		appointmentTime := startTime.Add(time.Duration(i) * interval)
		scheduleTimes = append(scheduleTimes, RoundTime(appointmentTime))
	}

	return scheduleTimes, nil
}
