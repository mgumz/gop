package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var (
	fversion bool
	vendor   string
	commitid string
	version  string

	rcfile        string
	shell         string
	verbose       bool
	gopath_marker string
)

// Usage function to helping the command line
var Usage = func(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n\n", msg)
	}
	fmt.Fprintf(os.Stderr, "Usage of %s\n\n    golo command [arguments]%v\n", path.Base(os.Args[0]), doc)
	flag.PrintDefaults()
	os.Exit(1)
}

var Version = func(msg string) {
	doc :=
		`golo
Version: %s
Commit ID: %s
`
	os_exit(1, os.Stderr, fmt.Sprintf(doc, version, commitid))
}

func os_exit(code int, writer io.Writer, msg string) {
	if msg != "" {
		fmt.Fprintln(writer, msg)
	}
	os.Exit(1)
}

func exists(dir, base string) (found bool, err error) {
	if base == "" {
		return false, fmt.Errorf("skip")
	}
	if _, err = os.Stat(filepath.Join(dir, base)); err != nil {
		return false, err
	}
	return true, nil
}

func main() {
	var gopath, vendorpath string

	flag.Parse()

	if fversion {
		Version("")
	}

	if rcfile != "" {
		if err := read_flagfile(flag.CommandLine, rcfile, os.Args[1:]); err != nil {
			// a missing ~/.golorc is not an error
			if !(os.IsNotExist(err) && rcfile == default_flagfile()) {
				os_exit(1, os.Stderr, fmt.Sprintf("Error parsing %q; %v\n", rcfile, err))
			}
		}
		if verbose {
			fmt.Printf("golo used flag file: %v\n", rcfile)
		}
	}

	// start looking for a '.gopath' file
	dir, err := os.Getwd()
	if err != nil {
		os_exit(2, os.Stderr, fmt.Sprintf("Error getting current dir %v\n", err))
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		os_exit(2, os.Stderr, fmt.Sprintf("Error making %q absolute; %v\n", dir, err))
	}

	// filepath.Clean(): The returned path ends in a slash only if it
	// represents a root directory, such as "/" on Unix or `C:\` on Windows.
	dir = filepath.Clean(dir)
	for dir[len(dir)-1] != filepath.Separator {
		if _, err := exists(dir, vendor); vendorpath == "" && err == nil {
			vendorpath = filepath.Join(dir, vendor)
			if verbose {
				fmt.Printf("golo picked vendorpath: %q\n", vendorpath)
			}
		} else if _, err := exists(dir, gopath_marker); gopath == "" && err == nil {
			gopath = dir
			if verbose {
				fmt.Printf("golo picked gopath: %q, found %q\n", gopath, filepath.Join(gopath, gopath_marker))
			}
		}
		dir = filepath.Clean(filepath.Dir(dir))
	}

	envpath := ""
	if vendorpath != "" {
		envpath += string(filepath.ListSeparator) + vendorpath
	}
	if gopath != "" {
		envpath += string(filepath.ListSeparator) + gopath
	}
	if envpath != "" {
		envpath = envpath[1:]
		os.Setenv("GOLO_MARKER", filepath.Join(envpath, gopath_marker))
		os.Setenv("GOPATH", envpath)
	}

	// golo command [arguments]
	cmdargs := strings.Join(flag.Args(), " ")
	switch flag.Arg(0) {
	case "build", "clean", "env", "fix", "fmt", "generate", "get", "help", "install", "list", "run", "test", "tool", "version", "vet":
		cmdargs = "go " + cmdargs
	}

	var cflag string
	shell, cflag = get_shell(shell)
	cmd := exec.Command(shell, cflag, cmdargs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if verbose {
		fmt.Printf("golo cmd.path: %q\n", cmd.Path)
		fmt.Printf("golo cmd.args: %v\n", cmd.Args)
		fmt.Printf("golo cmd.dir: %q\n", cmd.Dir)
		for i := range cmd.Env {
			fmt.Printf("golo cmd.env: %s\n", cmd.Env[i])
		}
	}

	if err = cmd.Run(); err != nil {
		os_exit(3, os.Stderr, fmt.Sprintf("Error statring the command; %v\n", err))
	}
}

func init() {
	flag.BoolVar(&fversion, "version", false, "display version info and exit")
	flag.StringVar(&vendor, "vendor", "", "look for vendor folder too")

	flag.StringVar(&rcfile, "rc", default_flagfile(), "read flags from given file")
	flag.StringVar(&shell, "shell", "", "shell to use to spawn 'go'")
	flag.BoolVar(&verbose, "verbose", false, "be more verbose")
	flag.StringVar(&gopath_marker, "gopath", ".gopath", "look for file marking the GOPATH folder")
}
