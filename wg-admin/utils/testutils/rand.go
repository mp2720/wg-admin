package testutils

import (
	"math/rand/v2"
	"strings"
)

func RandRangeInt64(from int64, to int64) int64 {
	return from + rand.Int64N(to-from+1)
}

func RandBool() bool {
	return RandRangeInt64(0, 1) == 0
}

// if runes is nil, then generated string can have any unicode character
func RandString(minLen int, maxLen int, runes []rune) string {
	strBuilder := strings.Builder{}
	strLen := RandRangeInt64(int64(minLen), int64(maxLen))
	for range strLen {
		var r rune
		if runes != nil {
			r = runes[rand.IntN(len(runes))]
		} else {
			r = rand.Int32()
		}
		strBuilder.WriteRune(r)
	}
	return strBuilder.String()
}

func RandValueOrNil[T any](v T) *T {
	if RandBool() {
		return &v
	} else {
		return nil
	}
}
