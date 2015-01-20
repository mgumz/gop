// +build !windows

package main

import "os"

const _DEFAULT_SHELL = "/bin/sh"
const _DEFAULT_CMD_FLAG = "-c"

func get_shell() (string, string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = _DEFAULT_SHELL
	}
	return shell, _DEFAULT_CMD_FLAG
}
