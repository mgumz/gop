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
	commands = ",build,clean,env,fix,fmt,generate,get,help,install,list,run,test,tool,version,vet,"
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

func main() {
	var err error
	flag.Parse()
	if help {
		Usage("")
	}
	// start looking for a '.gopath' file
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current dir %v\n", err)
		return
	}
	found := false
	for dir != "/" {
		if _, err := os.Stat(dir + "/.gopath"); err == nil {
			found = true
			break
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error getting current dir %v\n", err)
			return
		}
		dir = path.Dir(dir)
		err = os.Chdir(dir)
	}
	if found {
		os.Setenv("GOPATH", dir)
	}

	// golo command [arguments]
	cmdargs := strings.Join(flag.Args(), " ")
	if strings.Index(commands, ","+flag.Arg(0)+",") >= 0 {
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
	flag.BoolVar(&help, "help", false, "display this help scre")
}
