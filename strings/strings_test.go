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
