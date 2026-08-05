package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
	"time"
)

func run(name string) error {
	c, err := newConfig(name)
	if err != nil {
		return err
	}
	conn, err := ft12.NewConn(c.Device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	handlers := make(baos.Handlers)
	for _, sw := range c.Switches {
		handlers.Bool([]baos.DatapointID{sw}, func(id baos.DatapointID, value bool) error {
			_, _ = fmt.Fprintf(os.Stdout, "switch[datapoint=%d,value=%t]\n", id, value)
			return nil
		})
	}
	for _, th := range c.Thermostats {
		handlers.Float([]baos.DatapointID{th}, func(id baos.DatapointID, value float64) error {
			_, _ = fmt.Fprintf(os.Stdout, "thermostat[datapoint=%d,value=%0.1f]\n", id, value)
			return nil
		})
	}
	return runWith(conn, handlers)
}

func runWith(conn *ft12.Conn, handlers baos.Handlers) error {
	worker := baos.NewWorker(conn, baos.WithHandlers(handlers))
	defer worker.Stop()
	if err := worker.Run(); err != nil {
		return err
	}
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: events -config <config file>\n")
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.PrintDefaults()
}

func main() {

	var conf string
	flag.StringVar(&conf, "config", "", "config file to read")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() || conf == "" {
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := run(conf); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
