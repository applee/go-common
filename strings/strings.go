// Package strings provides some utility function on string.
package strings

import (
	"reflect"
	"unicode/utf8"
	"unsafe"
)

func ToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func FromBytes(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func SubString(s string, start int, end int) string {
	var i, j, idx int
	if start > end {
		return s
	}
	if start < 0 || end < 0 {
		l := utf8.RuneCountInString(s)
		if start < 0 {
			start = l + start
		}
		if end < 0 {
			end = l + end
		}
	}
	for j = range s {
		if i == start {
			idx = j
		}
		if i == end && end > 0 {
			return s[idx:j]
		}
		i++
	}
	return s[idx:]
}

func SubStr(s string, start int, length int) string {
	var i, j, idx int
	if length <= 0 {
		return s
	}
	if start < 0 {
		start = utf8.RuneCountInString(s) + start
	}
	for j = range s {
		if i == start {
			idx = j
		}
		if i == start+length {
			return s[idx:j]
		}
		i++
	}
	return s[idx:]
}
