package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	help     bool
	fversion bool
	vendor   string
	commitid string
	version  string
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
	fmt.Fprintf(os.Stderr, doc, version, commitid)
	os.Exit(1)
}

func exists(dir, base string) (found bool, err error) {
	if base == "" {
		return false, fmt.Errorf("skip")
	}
	if _, err = os.Stat(dir + "/" + base); err != nil {
		return false, err
	}
	return true, nil
}

func main() {
	var err error
	flag.Parse()
	if help {
		Usage("")
	}
	if fversion {
		Version("")
	}
	// start looking for a '.gopath' file
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current dir %v\n", err)
		return
	}
	gopath := ""
	vendorpath := ""
	for dir != "/" {
		if _, err := exists(dir, vendor); vendorpath == "" && err == nil {
			vendorpath = dir + "/" + vendor
		} else if _, err := exists(dir, ".gopath"); gopath == "" && err == nil {
			gopath = dir
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error getting current dir %v\n", err)
			return
		}
		dir = path.Dir(dir)
		err = os.Chdir(dir)
	}
	envpath := ""
	if vendorpath != "" {
		envpath += ":" + vendorpath
	}
	if gopath != "" {
		envpath += ":" + gopath
	}
	if envpath != "" {
		envpath = envpath[1:]
		os.Setenv("GOPATH", envpath)
	}

	// golo command [arguments]
	cmdargs := strings.Join(flag.Args(), " ")
	switch flag.Arg(0) {
	default:
	case "build", "clean", "env", "fix", "fmt", "generate", "get", "help", "install", "list", "run", "test", "tool", "version", "vet":
		cmdargs = "go " + cmdargs
	}

	cmd := exec.Command("/bin/bash", "-c", cmdargs)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating pipe %v\n", err)
		return
	}
	if err = cmd.Start(); err != nil {
		fmt.Printf("Error starting the command; %v\n", err)
		return
	}
	r := bufio.NewReader(stdout)
	for err == nil {
		if line, _, err := r.ReadLine(); err == nil {
			fmt.Println(string(line))
			continue
		}
		break
	}
	cmd.Wait()
}

func init() {
	flag.BoolVar(&help, "help", false, "display this help screeen")
	flag.BoolVar(&fversion, "version", false, "display version info and exit")
	flag.StringVar(&vendor, "vendor", "", "look for vendor folder too")
}
