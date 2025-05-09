package utils

import (
	"errors"
	"time"
)

const (
	minimumAcceptableTimesPerDay = 1
	maximumAcceptableTimesPerDay = 24
	dayStartHour                 = 8
	dayEndHour                   = 22
)

func GenerateScheduleTimes(day time.Time, timesPerDay uint64, nowFunc func() time.Time) ([]time.Time, error) {
	if timesPerDay < minimumAcceptableTimesPerDay || timesPerDay > maximumAcceptableTimesPerDay {
		return nil, errors.New("неверное количество приемов в день")
	}

	var scheduleTimes []time.Time

	createdService := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(day.Year(), day.Month(), day.Day(), dayStartHour, 0, 0, 0, day.Location())
	endTime := time.Date(day.Year(), day.Month(), day.Day(), dayEndHour, 0, 0, 0, day.Location())
	duration := endTime.Sub(startTime)
	now := nowFunc()

	if day.Year() < createdService.Year() || day.Year() > now.Year() {
		return nil, errors.New("неверный год")
	}

	if day.IsZero() {
		return nil, errors.New("некорректная дата: zero time")
	}

	if now.After(endTime) {
		day = day.AddDate(0, 0, 1)
		startTime = time.Date(day.Year(), day.Month(), day.Day(), dayStartHour, 0, 0, 0, day.Location())
	}

	if timesPerDay == minimumAcceptableTimesPerDay {
		scheduleTimes = append(scheduleTimes, startTime.Add(duration/2))
		return scheduleTimes, nil
	}

	if timesPerDay == maximumAcceptableTimesPerDay {
		for i := uint64(0); i <= 14; i++ {
			scheduleTimes = append(scheduleTimes, startTime.Add(time.Duration(i)*time.Hour))
		}
		return scheduleTimes, nil
	}

	interval := duration / time.Duration(timesPerDay-1)

	for i := uint64(0); i < timesPerDay; i++ {
		appointmentTime := startTime.Add(time.Duration(i) * interval)
		scheduleTimes = append(scheduleTimes, RoundTime(appointmentTime))
	}

	return scheduleTimes, nil
}
