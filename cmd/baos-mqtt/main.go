package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"thermo-watchdog/internal/config"
)

func halt(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "Usage: thermo-watchdog -config <config file>\n")
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.PrintDefaults()
}

func main() {
	var conf string
	flag.StringVar(&conf, "config", "", "path to config file")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() || conf == "" {
		usage(os.Stderr)
		os.Exit(2)
	}

	c, err := config.Parse(conf)
	if err != nil {
		halt(err)
	}
	if err := run(c); err != nil {
		halt(err)
	}
}
