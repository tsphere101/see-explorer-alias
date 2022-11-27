package wt

import (
	"os/exec"
)

// run windows terminal with the specified path
func RunWt(path string, program string) {
	// Run windows terminal in new process

	cmd := exec.Command(program, "-d", path)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

}
