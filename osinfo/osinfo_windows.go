package osinfo

import "strings"

func Gather() (i *OSInfo) {
	out, err := CmdOut("wmic", "os", "get", "Caption,OSArchitecture,Version", "/value")
	for err != nil {
		return
	}
	info := strings.Split(out, "\n")
	i = &OSInfo{Kernel: strings.Split(info[0], "=")[1],
		Release:  strings.Split(info[2], "=")[1],
		Platform: strings.Split(info[1], "=")[1]}
	return

}
