package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func default_flagfile() string {
	home, _ := homedir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "_golorc")
	}
	return filepath.Join(home, ".golorc")
}

// read file 'path', line by line. the lines are used to
// override any flags set in 'to'. if no errors occured, 'flags'
// is used to override the flags given in 'path'.
func read_flagfile(to *flag.FlagSet, path string, flags []string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// read in the file line by line
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return err
	}

	// override original values by lines
	if err := to.Parse(lines); err != nil {
		return err
	}

	// refill in original values
	return to.Parse(flags)
}

func homedir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Warning: can't detect current user: %v\n", err)
	}
	return user.HomeDir, nil
}
