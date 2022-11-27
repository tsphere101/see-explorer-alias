package vscode

import (
	"os/exec"
)

// Run vscode with the specified path
func RunCode(path string, program string) {
	cmd := exec.Command(program, path)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
