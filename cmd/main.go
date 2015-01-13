package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var help bool

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
	flag.Parse()
	if help {
		Usage("")
	}
	// golo mycommand custom args
}

func init() {
	flag.BoolVar(&help, "help", false, "display this help scre")
}
