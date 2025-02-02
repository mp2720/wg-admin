package utils

import (
	"time"
)

// Interface for mocking the time.Now function.
type Clock interface {
	Now() time.Time
}

// Just returns time.Now
type RealClock struct{}

func (rc RealClock) Now() time.Time {
	return time.Now()
}
