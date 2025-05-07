package tests

import (
	"dolittle2/internal/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestGenerateScheduleTimes_Negative1(t *testing.T) {
	invalidValues := []uint64{0, 25, 100, 123123124}
	for _, val := range invalidValues {
		t.Run(fmt.Sprintf("Invalid timesPerDay: %d", val), func(t *testing.T) {
			schedule, err := utils.GenerateScheduleTimes(time.Now(), val)
			assert.Error(t, err)
			assert.Nil(t, schedule)
		})
	}
}

func TestGenerateScheduleTimes_Negative2(t *testing.T) {
	originalNowFunc := utils.NowFunc
	defer func() { utils.NowFunc = originalNowFunc }()

	mockNow := time.Date(2025, 5, 7, 23, 0, 0, 0, time.UTC)
	utils.NowFunc = func() time.Time { return mockNow }

	day := time.Date(2025, 5, 7, 0, 0, 0, 0, time.UTC)
	schedule, err := utils.GenerateScheduleTimes(day, 3)

	assert.NoError(t, err)
	assert.NotNil(t, schedule)
	assert.Len(t, schedule, 3)

	expectedDay := day.AddDate(0, 0, 1).Day()
	for _, appointment := range schedule {
		assert.Equal(t, expectedDay, appointment.Day(), "Appointment should be on the next day")
		slog.Info("Test appointment day", slog.Int("day", appointment.Day()))
	}
}

func TestGenerateScheduleTimes_Negative3(t *testing.T) {
	zeroTime := time.Time{}
	schedule, err := utils.GenerateScheduleTimes(zeroTime, 3)
	assert.Error(t, err)
	assert.Nil(t, schedule)
	slog.Info("sda", slog.Any("day", schedule))
}

func TestGenerateScheduleTimes_Negative4(t *testing.T) {
	test := []time.Time{time.Date(1000, 5, 7, 23, 0, 0, 0, time.UTC),
		time.Date(5000, 5, 7, 23, 0, 0, 0, time.UTC)}
	for _, tt := range test {
		t.Run(fmt.Sprintf("Test old day: %d", tt.Day()), func(t *testing.T) {
			schedule, err := utils.GenerateScheduleTimes(tt, 1)
			assert.Error(t, err)
			assert.Nil(t, schedule)
		})
	}
}

func TestGenerateScheduleTimes_Negative5(t *testing.T) {
	tpd := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	for _, tt := range tpd {
		t.Run(fmt.Sprintf("Тест на кратность 15: %d", tt), func(t *testing.T) {
			schedule, err := utils.GenerateScheduleTimes(time.Now(), tt)
			assert.NoError(t, err)
			assert.NotNil(t, schedule)

			expectedLen := int(tt)
			if tt == 24 {
				expectedLen = 15
			}

			assert.Len(t, schedule, expectedLen)

			for _, s := range schedule {
				minutes := s.Minute()
				if minutes%15 != 0 {
					t.Errorf("Время %v не кратно 15 минутам", s.Format("15:04"))
				}
			}
		})
	}
}
