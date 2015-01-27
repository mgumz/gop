package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	NAME   = "gop"
	MARKER = ".gopath"
)

var (
	_VERSION_MAJOR = "1"
	_VERSION_MINOR = "x"
	_DATE          = "unknown"
	_COMMIT_ID     = "unknown"
)

func main() {
	var (
		gopath     string
		vendorpath string
		cli        = struct {
			rcfile       string
			shell        string
			marker       string
			vendor       string
			verbose      bool
			show_version bool
		}{marker: MARKER}
	)
	flag.BoolVar(&cli.show_version, "version", cli.show_version, "display version info and exit")
	flag.StringVar(&cli.vendor, "vendor", cli.vendor, "look for vendor folder too")
	flag.StringVar(&cli.rcfile, "rc", default_flagfile(), "read flags from given file")
	flag.StringVar(&cli.shell, "shell", cli.shell, "shell to use to spawn 'go'")
	flag.BoolVar(&cli.verbose, "verbose", cli.verbose, "be more verbose")
	flag.StringVar(&cli.marker, "gopath", cli.marker, "look for file marking the GOPATH folder")

	flag.Parse()

	if cli.show_version {
		print_version()
		os.Exit(0)
	}

	if cli.rcfile != "" {
		if err := read_flagfile(flag.CommandLine, cli.rcfile, os.Args[1:]); err != nil {
			// a missing default_flagfile() is not an error
			if !(os.IsNotExist(err) && cli.rcfile == default_flagfile()) {
				os_exit(1, os.Stderr, fmt.Sprintf("error parsing %q; %v\n", cli.rcfile, err))
			}
		}
		if cli.verbose {
			fmt.Printf("%s used flag file: %v\n", NAME, cli.rcfile)
		}
	}

	// start looking for the marker filec
	dir, err := os.Getwd()
	if err != nil {
		os_exit(2, os.Stderr, fmt.Sprintf("error getting current dir %v\n", err))
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		os_exit(2, os.Stderr, fmt.Sprintf("error making %q absolute; %v\n", dir, err))
	}

	// filepath.Clean(): The returned path ends in a slash only if it
	// represents a root directory, such as "/" on Unix or `C:\` on Windows.
	dir = filepath.Clean(dir)
	for dir[len(dir)-1] != filepath.Separator {
		if _, err := exists(dir, cli.vendor); vendorpath == "" && err == nil {
			vendorpath = filepath.Join(dir, cli.vendor)
			if cli.verbose {
				fmt.Printf("%s picked vendorpath: %q\n", NAME, vendorpath)
			}
		} else if _, err := exists(dir, cli.marker); gopath == "" && err == nil {
			gopath = dir
			if cli.verbose {
				fmt.Printf("%s picked gopath: %q, found %q\n", NAME, gopath, filepath.Join(gopath, cli.marker))
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
		os.Setenv("GOP_MARKER", filepath.Join(envpath, cli.marker))
		os.Setenv("GOPATH", envpath)
	}

	// golo command [arguments]
	cmdargs := strings.Join(flag.Args(), " ")
	switch flag.Arg(0) {
	case "build", "clean", "env", "fix", "fmt", "generate", "get",
		"help", "install", "list", "run", "test", "tool",
		"version", "vet":
		cmdargs = "go " + cmdargs
	}

	var cflag string
	cli.shell, cflag = get_shell(cli.shell)
	cmd := exec.Command(cli.shell, cflag, cmdargs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if cli.verbose {
		fmt.Printf("%s cmd.path: %q\n", NAME, cmd.Path)
		fmt.Printf("%s cmd.args: %v\n", NAME, cmd.Args)
		fmt.Printf("%s cmd.dir: %q\n", NAME, cmd.Dir)
		for i := range cmd.Env {
			fmt.Printf("%s cmd.env: %s\n", NAME, cmd.Env[i])
		}
	}

	if err = cmd.Run(); err != nil {
		os_exit(3, os.Stderr, fmt.Sprintf("error running the command; %v\n", err))
	}
}

// Usage function to helping the command line
func usage(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", msg)
	}
	fmt.Fprintf(os.Stderr, "usage of %s\n\n    %s command [arguments]%v\n", path.Base(os.Args[0]), NAME, doc)
	flag.PrintDefaults()
	os.Exit(0)
}

func print_version() {
	fmt.Fprintf(os.Stderr,
		"%s\nversion: %s.%s\ngit-commit: %s\ndate: %s\ngo: %s\n\n",
		NAME,
		_VERSION_MAJOR, _VERSION_MINOR,
		_COMMIT_ID,
		_DATE,
		runtime.GOOS+"-"+runtime.Version())
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
