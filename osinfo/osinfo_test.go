package osinfo

import "testing"

func TestGetInfo(t *testing.T) {
	info := Gather()
	if info == nil {
		t.Fatal("gather os info failed.")
	}
	t.Log(info)
}
