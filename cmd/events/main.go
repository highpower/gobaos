package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
)

func run(name string) error {
	c, err := newConfig(name)
	if err != nil {
		return err
	}
	dm := datapointMap{}
	for _, d := range c.Datapoints {
		for _, id := range d.IDs {
			if err := dm.add(id, d.Type); err != nil {
				return err
			}
		}
	}
	conn, err := ft12.NewConn(c.Device)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, dm)
}

func runWith(conn *ft12.Conn, dm datapointMap) error {
	for {
		if err := conn.Poll(); err != nil {
			return err
		}
		buffer := make([]byte, conn.Available())
		if _, err := conn.Read(buffer); err != nil {
			return err
		}
		if !baos.IsDatapointValueInd(buffer) {
			continue
		}
		indication := baos.DatapointValueInd{}
		if err := indication.UnmarshalBinary(buffer); err != nil {
			return err
		}
		if err := dm.handleIndication(&indication); err != nil {
			return err
		}
	}
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
