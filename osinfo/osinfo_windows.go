package osinfo

import "strings"

func Gather() (i *OSInfo) {
	out, err := CmdOut("wmic", "os", "get", "Caption,Version,OSArchitecture", "/value")
	for err != nil {
		return
	}
	info := strings.Split(out, " ")
	i = &OSInfo{Kernel: info[0], Release: info[1], Platform: info[2]}
	return

}
