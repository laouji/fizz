package fizzbuzz

import "fmt"

func Calculate(int1, int2 int, str1, str2 string, limit int) []string {
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
	return out
}
