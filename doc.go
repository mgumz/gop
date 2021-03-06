package main

const doc = `

The commands are:

build       compile packages and dependencies
clean       remove object files
env         print Go environment information
fix         run go tool fix on packages
fmt         run gofmt on package sources
generate    generate Go files by processing source
get         download and install packages and dependencies
install     compile and install packages and dependencies
list        list packages
run         compile and run Go program
test        test packages
tool        run specified go tool
version     print Go version
vet         run go tool vet on packages

Use "gop help [command]" for more information about a command.

Additional help topics:

c           calling between Go and C
filetype    file types
gopath      GOPATH environment variable
importpath  import path syntax
packages    description of package lists
testflag    description of testing flags
testfunc    description of testing functions

Use "gop help [topic]" for more information about that topic.
`
