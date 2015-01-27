// +build windows

package main

import (
	"os"
	"path/filepath"
)

const _DEFAULT_SHELL = "cmd.exe"
const _DEFAULT_CMD_FLAG = "/c"

func get_shell(shell string) (string, string) {

	// TODO: detect that we run inside powershell

	if shell == "" {
		shell = os.Getenv("COMSPEC")
	}
	if shell == "" {
		shell = _DEFAULT_SHELL
	} else if filepath.Base(shell) == "powershell.exe" {
		return shell, "-Command"
	}

	return shell, _DEFAULT_CMD_FLAG
}
