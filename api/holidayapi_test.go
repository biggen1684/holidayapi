package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekday(t *testing.T) {
	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "Thursday", weekday(date))
}

func TestUnderThirty_False(t *testing.T) {
	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.False(t, underThirty(date))
}

func TestUnderThirty_True(t *testing.T) {
	date := time.Now().Add(24 * time.Hour)
	assert.True(t, underThirty(date))
}
