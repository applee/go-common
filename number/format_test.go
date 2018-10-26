package number

import (
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		num int64
		fmt string
	}{
		{123, "123"},
		{1234, "1,234"},
		{123456, "123,456"},
	}
	for i := range cases {
		got := Format(cases[i].num, ",")
		if got != cases[i].fmt {
			t.Fatalf("%d after Round want %s, got %s",
				cases[i].num, cases[i].fmt, got)
		}
	}
}
