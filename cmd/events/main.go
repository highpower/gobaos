package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
)

func handleDatapointValueInd(buffer []byte) error {
	indication := baos.DatapointValueInd{}
	if err := indication.UnmarshalBinary(buffer); err != nil {
		return err
	}
	switch indication.StartDatapoint {
	case 2, 3, 4, 5, 6, 7:
		dpt, err := indication.Values[0].DPTFloat()
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: %s\n", indication.Values[0].ID, dpt.String())
	case 38, 39, 40:
		dpt, err := indication.Values[0].DPTBool()
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: %s\n", indication.Values[0].ID, dpt.String())
	default:
		_, _ = fmt.Fprintf(os.Stdout, "happen value indication of %d datapoint\n", indication.StartDatapoint)
	}
	return nil
}

func run(device string) error {
	conn, err := ft12.NewConn(device)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	for {
		if _, err := conn.Read(nil); err != nil {
			return err
		}
		l := conn.Available()
		buffer := make([]byte, l)
		_, err := conn.Read(buffer)
		if err != nil {
			return err
		}
		if baos.IsDatapointValueInd(buffer) {
			if err := handleDatapointValueInd(buffer); err != nil {
				return err
			}
		}
	}
}

func usage(writer io.Writer) {
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.PrintDefaults()
}

func main() {

	var device string
	flag.StringVar(&device, "device", "", "device to read uart")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() || device == "" {
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := run(device); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
