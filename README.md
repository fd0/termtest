[![GoDoc](https://godoc.org/github.com/fd0/termtest?status.svg)](http://godoc.org/github.com/fd0/termtest)
[![Build Status](https://travis-ci.org/fd0/termtest.svg?branch=master)](https://travis-ci.org/fd0/termtest)

termtest is a small helper program to run terminal user interface tests by
using tmux in Go.

Example:

```go
// create terminal
term, _ := termtest.New()

buf, _ := term.Run(12, 13, "echo This is a long line that will wrap for sure")

// print the output string:
// Output: "This is a lo\nng line that\n will wrap f\nor sure\n\n\n\n\n\n\n\n\n\n"
fmt.Printf("%q\n", buf)

// stop the tmux instance
term.Exit()
```
