package slice

import (
	"errors"
	"reflect"
)

// DiffString returns the subtraction of tow string slices
func DiffString(slice1, slice2 []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func RemoveDuplicatesDisorder() {

}

// RemoveIntUnordered returns the unordered remaining slice after getting rid
// of arg2 from arg1. This function do not preserve the order of origin slice.
func RemoveIntUnordered(data, remove []int) []int {
	n := len(data)
	i := 0
loop:
	for i < n {
		d := data[i]
		for _, r := range remove {
			if r == d {
				data[i] = data[n-1]
				n--
				continue loop
			}
		}
		i++
	}
	return data[0:n]
}

// RemoveStringUnordered returns the unordered remaining slice after getting
// rid of arg2 from arg1.
func RemoveStringUnordered(data, remove []string) []string {
	n := len(data)
	i := 0
loop:
	for i < n {
		d := data[i]
		for _, r := range remove {
			if r == d {
				data[i] = data[n-1]
				n--
				continue loop
			}
		}
		i++
	}
	return data[0:n]
}

// RemoveUnordered returns the unordered remaining slice after getting rid of
// arg2 from arg1 using reflection. Please only use it in generic type
// processing for low performance.
func RemoveUnordered(data, remove interface{}) (interface{}, error) {
	d := reflect.ValueOf(data)
	r := reflect.ValueOf(remove)
	if d.Type().Kind() != reflect.Slice || r.Type().Kind() != reflect.Slice {
		return nil, errors.New("Invalid type.")
	}

	var (
		di, ri reflect.Value
		i, j   int
		m      = d.Len()
		n      = r.Len()
	)
loop:
	for i < m {
		di = d.Index(i)
		for j = 0; j < n; j++ {
			ri = r.Index(j)
			if ri.Interface() == di.Interface() {
				d.Index(i).Set(d.Index(m - 1))
				m--
				continue loop
			}
		}
		i++
	}
	return d.Slice(0, m).Interface(), nil
}
