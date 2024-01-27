package fizzbuzz_test

import (
	"testing"

	"github.com/laouji/fizz/pkg/fizzbuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type fizzBuzzTestSuite struct {
	suite.Suite
}

func TestFizzBuzz(t *testing.T) {
	suite.Run(t, new(fizzBuzzTestSuite))
}

func (s *fizzBuzzTestSuite) TestCalculate() {
	cases := map[string]struct {
		int1   int
		int2   int
		str1   string
		str2   string
		limit  int
		result []string
		field  string
		err    error
	}{
		"classic fizzbuzz": {
			int1:   3,
			int2:   5,
			str1:   "fizz",
			str2:   "buzz",
			limit:  16,
			result: []string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16"},
		},
		"alternate strings": {
			int1:   2,
			int2:   6,
			str1:   "two",
			str2:   "six",
			limit:  10,
			result: []string{"1", "two", "3", "two", "5", "twosix", "7", "two", "9", "two"},
		},
		"limit out of range": {
			int1:  0,
			int2:  1,
			limit: fizzbuzz.CalculateMaxLimit + 1,
			field: "limit",
			err:   fizzbuzz.ErrOutOfRange,
		},
		"int1 is 0": {
			int1:  0,
			limit: fizzbuzz.CalculateMaxLimit,
			field: "int1",
			err:   fizzbuzz.ErrOutOfRange,
		},
		"int2 is too low": {
			int1:  4,
			int2:  -4,
			limit: fizzbuzz.CalculateMaxLimit,
			field: "int2",
			err:   fizzbuzz.ErrOutOfRange,
		},
	}

	for testName, c := range cases {
		s.T().Run(testName, func(t *testing.T) {
			res, field, err := fizzbuzz.Calculate(c.int1, c.int2, c.str1, c.str2, c.limit)
			require.Equal(t, c.err, err)
			assert.Equal(t, c.result, res)
			assert.Equal(t, c.field, field)
			if c.err == nil {
				assert.Len(t, res, c.limit)
			}
		})
	}
}
