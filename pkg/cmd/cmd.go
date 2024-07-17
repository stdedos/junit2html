package cmd

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/stdedos/junit2html/pkg/logging"

	"github.com/stdedos/junit2html/pkg/convert"
	"github.com/stdedos/junit2html/pkg/parse"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	Quiet   []bool `short:"q" long:"quiet" description:"Decrease verbosity"`
}

var opts Options

func EntryPoint(args []string) (string, error) {
	positionalArgs, err := flags.ParseArgs(&opts, args)

	if flags.WroteHelp(err) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("error parsing flags: %w", err)
	}

	if opts.Verbose != nil || opts.Quiet != nil {
		by := len(opts.Quiet) - len(opts.Verbose)

		if by != 0 {
			logging.ModifyVerbosity(by)
		}
	}

	files := parse.Files(positionalArgs)
	suites := parse.Suites(files)

	return convert.Convert(suites, files)
}
