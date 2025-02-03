package testutils

import "time"

const (
	firstRandTimeMin = 1000000000
	firstRandTimeMax = 4000000000

	randTimeChangeMin = 1
	randTimeChangeMax = 100000

	maxNanoseconds = 999999999
)

// Generate random time values for tests.
type RandomClock struct {
	NextReturnValue time.Time
	DoAscend        bool
}

func NewRandomClock(doAscend bool) RandomClock {
	secs := RandRangeInt64(firstRandTimeMin, firstRandTimeMax)
	nsecs := RandRangeInt64(0, maxNanoseconds)

	c := RandomClock{
		DoAscend:        doAscend,
		NextReturnValue: time.Unix(secs, nsecs),
	}

	return c
}

func (c *RandomClock) Now() time.Time {
	var secs int64
	if c.DoAscend {
		secs = c.NextReturnValue.Unix() +
			RandRangeInt64(randTimeChangeMin, randTimeChangeMax)
	} else {
		secs = RandRangeInt64(firstRandTimeMin, firstRandTimeMax)
	}
	nsecs := c.NextReturnValue.Unix()

	ret := c.NextReturnValue

	c.NextReturnValue = time.Unix(secs, nsecs)

	return ret
}
