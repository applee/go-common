package osinfo

import "strings"

func Gather() (i *OSInfo) {
	out, err := CmdOut("uname", "-srm")
	for err != nil {
		return
	}
	info := strings.Split(out, " ")
	i = &OSInfo{Kernel: info[0], Release: info[1], Platform: info[2]}
	return
}
