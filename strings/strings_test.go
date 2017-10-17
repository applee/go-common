package strings

import "testing"

func BenchmarkS2B(b *testing.B) {
	var result []byte
	str := "abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg"
	for i := 0; i < b.N; i++ {
		result = []byte(str)
	}
	if len(result) == 0 {
		b.Error("failed")
	}
}

func BenchmarkS2BFast(b *testing.B) {
	var result []byte
	str := "abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg"
	for i := 0; i < b.N; i++ {
		result = ToBytes(str)
	}
	if len(result) == 0 {
		b.Error("failed")
	}
}

func TestSubString(t *testing.T) {
	s := "你好hello中国"
	cases := []struct {
		start, end int
		want       string
	}{
		{0, 2, "你好"},
		{-2, -1, "中"},
		{-1, 0, "国"},
	}

	for idx, c := range cases {
		if got := SubString(s, c.start, c.end); got != c.want {
			t.Fatalf("case %d: want: %s, got: %s", idx, c.want, got)
			break
		}
	}
}
