gop
====

*gop* looks recursively for a folder marked by a file (default: ".gopath")
and sets Go's GOPATH environment to that folder. It then executes the
the given command.

*gop* is a fork of [golo][golo] by Richard Bucker (which itself got inspired
by [goli][goli], written by James Lawrence)

Usage
-----

    $> cd ~/my_workspace
    $> touch .gopath
    $> cd src/deep/into/my/project
    $> gop -verbose build -v

LICENSE
-------

MIT


[golo]: https://bitbucket.org/oneoffcode/golo
[goli]: https://bitbucket.org/jatone/gilo
