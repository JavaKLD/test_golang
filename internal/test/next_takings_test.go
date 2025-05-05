package test

import (
	"dolittle2/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateScheduleTimes_ValidTimesPerDay(t *testing.T) {
	day := time.Date(2025, 5, 3, 0, 0, 0, 0, time.UTC)
	timesPerDay := uint64(3)
	scheduleTimes, err := utils.GenerateScheduleTimes(day, timesPerDay)

	assert.NoError(t, err)

	assert.Equal(t, timesPerDay, uint64(len(scheduleTimes)))

	assert.True(t, scheduleTimes[0].After(day))
	assert.True(t, scheduleTimes[len(scheduleTimes)-1].Before(day.Add(24*time.Hour)))
}
