package slice

import "testing"
import "reflect"
import "sort"

func TestRemoveDisorder(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5"}
	b := []string{"1", "2"}
	excepted := []string{"3", "4", "5"}
	c, err := RemoveUnordered(a, b)
	sort.Strings(c.([]string))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, excepted) {
		t.Fatal()
	}
}

func TestRemoveStringDisorder(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5"}
	b := []string{"1", "2"}
	excepted := []string{"3", "4", "5"}
	c := RemoveStringUnordered(a, b)
	sort.Strings(c)
	if !reflect.DeepEqual(c, excepted) {
		t.Fatal()
	}
}

func TestRemoveDuplicatesInt(t *testing.T) {
	s := []int{1, 2, 3, 1, 2, 3}
	excepted := []int{1, 2, 3}
	RemoveDuplicatesInt(&s)
	if !reflect.DeepEqual(s, excepted) {
		t.Fatal()
	}
}

func TestContainsInt(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	a := 3
	b := 8
	if !ContainsInt(a, s) {
		t.Fatal()
	}
	if ContainsInt(b, s) {
		t.Fatal()
	}

}

func BenchmarkRemoveStringDisorder(b *testing.B) {
	m := []string{"1", "2", "3", "4", "5"}
	n := []string{"1", "2"}
	for i := 0; i < b.N; i++ {
		RemoveStringUnordered(m, n)
	}
}

func BenchmarkRemoveDisorder(b *testing.B) {
	m := []string{"1", "2", "3", "4", "5"}
	n := []string{"1", "2"}
	for i := 0; i < b.N; i++ {
		RemoveUnordered(m, n)
	}
}
