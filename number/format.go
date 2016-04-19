package number

import (
	"fmt"
	"math"
	"strconv"
)

// NumberFormat 千分位格式化数字
func NumberFormat(number float64, thousandSep string) string {
	number = Round(number)
	str := strconv.FormatInt(int64(number), 10)
	fmt.Println(str)
	nl := len(str)
	tl := len(thousandSep)
	rl := nl + (nl-1)/3*tl
	b := make([]byte, rl)
	if nl == rl {
		return str
	}

	fmt.Println(rl)
	count := 0
	for i, j := nl-1, rl-1; i >= 0; i, j = i-1, j-1 {
		b[j] = str[i]
		count++
		if count%3 == 0 && j >= 1 {
			copy(b[j-tl:j], thousandSep)
			j -= tl
		}
	}
	return string(b)
}

// Round 四舍五入
func Round(a float64) float64 {
	if a < 0 {
		return math.Ceil(a - 0.5)
	}
	return math.Floor(a + 0.5)
}
