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
