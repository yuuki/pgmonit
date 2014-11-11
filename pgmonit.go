package main

import (
	"os"
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
    Version bool `short:"v" long:"version" description:"Show version"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("version %s\n", Version)
		os.Exit(0)
	}

	fmt.Printf("args %s\n", args)
}
