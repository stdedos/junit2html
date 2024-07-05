package main

import (
	"os"

	"github.com/stdedos/junit2html/pkg/cmd"
)

func main() {
	cmd.EntryPoint(os.Args[1:])
}
