package tests

import (
	"dolittle2/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateScheduleTimes_Pozitive(t *testing.T) {
	tests := []struct {
		name           string
		day            time.Time
		timesPerDay    uint64
		expectedLength int
		expextedHours  []int
	}{
		{
			name:           "Один раз в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    1,
			expectedLength: 1,
			expextedHours:  []int{15},
		},
		{
			name:           "Два раза в денб",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    2,
			expectedLength: 2,
			expextedHours:  nil,
		},
		{
			name:           "Три раза в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    3,
			expectedLength: 3,
			expextedHours:  nil,
		},
		{
			name:           "Четыре раза в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    4,
			expectedLength: 4,
			expextedHours:  nil,
		},
		{
			name:           "Пять раз в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    5,
			expectedLength: 5,
			expextedHours:  nil,
		},
		{
			name:           "Шесть раз в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    6,
			expectedLength: 6,
			expextedHours:  nil,
		},
		{
			name:           "Восемь раз в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    8,
			expectedLength: 8,
			expextedHours:  nil,
		},
		{
			name:           "Двенадцать раз в день",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    12,
			expectedLength: 12,
			expextedHours:  nil,
		},
		{
			name:           "Каждый час (24)",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    24,
			expectedLength: 15,
			expextedHours:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := utils.GenerateScheduleTimes(tt.day, tt.timesPerDay, time.Now)
			assert.NoError(t, err)
			assert.Len(t, schedule, tt.expectedLength)

			for _, s := range schedule {
				assert.GreaterOrEqual(t, s.Hour(), 8)
				assert.LessOrEqual(t, s.Hour(), 22)
				assert.True(t, s.Minute()%15 == 0)
			}
			if len(tt.expextedHours) > 0 {
				for i, s := range schedule {
					assert.Equal(t, tt.expextedHours[i], s.Hour())
				}
			}

			if tt.expectedLength > 1 {
				diff := schedule[1].Sub(schedule[0])
				for i := 1; i < len(schedule); i++ {
					currentDiff := schedule[i].Sub(schedule[i-1])
					assert.LessOrEqual(t, absDuration(diff-currentDiff), 15*time.Minute, "Интервал отличается больше, чем на 15 минут")
				}
			}
			t.Log(tt.name, "тест пройден")
		})
	}
}
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
