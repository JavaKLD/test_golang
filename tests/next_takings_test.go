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
			name:           "Once per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    1,
			expectedLength: 1,
			expextedHours:  []int{8},
		},
		{
			name:           "Twice per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    2,
			expectedLength: 2,
			expextedHours:  nil,
		},
		{
			name:           "Three times per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    3,
			expectedLength: 3,
			expextedHours:  nil,
		},
		{
			name:           "Four times per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    4,
			expectedLength: 4,
			expextedHours:  nil,
		},
		{
			name:           "Five times per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    5,
			expectedLength: 5,
			expextedHours:  nil,
		},
		{
			name:           "Eight times per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    8,
			expectedLength: 8,
			expextedHours:  nil,
		},
		{
			name:           "Twelve times per day",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    12,
			expectedLength: 12,
			expextedHours:  nil,
		},
		{
			name:           "Every hour (24 times)",
			day:            time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			timesPerDay:    24,
			expectedLength: 24,
			expextedHours:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := utils.GenerateScheduleTimes(tt.day, tt.timesPerDay)
			assert.NoError(t, err)
			assert.Len(t, schedule, tt.expectedLength)

			for _, s := range schedule {
				hour := s.Hour()
				assert.GreaterOrEqual(t, hour, 8)
				assert.LessOrEqual(t, hour, 22)

				minute := s.Minute()
				assert.True(t, minute%15 == 0)
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
					assert.Equal(t, diff, currentDiff)
				}
			}
		})
	}
}
