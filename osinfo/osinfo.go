package osinfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type OSInfo struct {
	Kernel   string
	Release  string
	Platform string
}

func (this *OSInfo) String() string {
	return fmt.Sprintf("Kernel: %s, Release: %s, Platform: %s",
		this.Kernel, this.Release, this.Platform)
}

func CmdOut(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
