package number

import (
	"strconv"

	"github.com/applee/go-common/strings"
)

// Format formats the integer number with the thousands separator for display.
func Format(number int64, thousandSep string) string {
	str := strconv.FormatInt(number, 10)
	nl := len(str)
	tl := len(thousandSep)
	rl := nl + (nl-1)/3*tl
	b := make([]byte, rl)
	if nl == rl {
		return str
	}

	count := 0
	for i, j := nl-1, rl-1; i >= 0; i, j = i-1, j-1 {
		b[j] = str[i]
		count++
		if count%3 == 0 && j >= 1 {
			copy(b[j-tl:j], thousandSep)
			j -= tl
		}
	}
	return strings.FromBytes(b)
}
