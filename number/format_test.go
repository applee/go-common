package number

import (
	"testing"
)

func Test_FormatNumber(t *testing.T) {
	a := -10.5
	t.Log(Round(a))
	b := 11.2
	t.Log(Round(b))
	c := 1.8
	t.Log(Round(c))

	d := 123456.22
	t.Log(NumberFormat(d, ","))

	e := 12345678.567
	t.Log(NumberFormat(e, ","))
}
