package main

import (
	"flag"

	"github.com/stdedos/junit2html/pkg/cmd"
)

// arguments
var (
	xmlReports *string
)

func init() {
	xmlReports = flag.String("xmlReports", "", "Comma delimited glob expressions describing the files to scan")
}

func main() {
	flag.Parse()
	cmd.EntryPoint(*xmlReports)
}
