package fizzbuzz

import (
	"errors"
	"fmt"
)

const CalculateMaxLimit = 1999999

var ErrOutOfRange = errors.New("out of range")

func Calculate(int1, int2 int, str1, str2 string, limit int) ([]string, string, error) {
	if limit < 0 || limit > CalculateMaxLimit {
		return nil, "limit", ErrOutOfRange
	}
	if int1 < 1 || int1 > CalculateMaxLimit {
		return nil, "int1", ErrOutOfRange
	}
	if int2 < 1 || int2 > CalculateMaxLimit {
		return nil, "int2", ErrOutOfRange
	}

	out := make([]string, 0, limit)

	for i := 1; i <= limit; i++ {
		var val string
		if i%int1 == 0 {
			val += str1
		}

		if i%int2 == 0 {
			val += str2
		}

		if val == "" {
			val = fmt.Sprintf("%d", i)
		}

		out = append(out, val)
	}
	return out, "", nil
}
