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
	_, _ = fmt.Fprintf(os.Stdout, "happen value indication of %d datapoint\n", indication.StartDatapoint)
	return nil
}

func run(device string) {
	conn, err := ft12.NewConn(device)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "connect error: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	for {
		if _, err := conn.Read(nil); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read error: %s\n", err.Error())
			os.Exit(1)
		}
		l := conn.Available()
		buffer := make([]byte, l)
		_, err := conn.Read(buffer)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "buffer read error: %s\n", err.Error())
		}
		if baos.IsDatapointValueInd(buffer) {
			if err := handleDatapointValueInd(buffer); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "handle datapoint value error: %s\n", err.Error())
				os.Exit(1)
			}
		}
	}
}

func usage(writer io.Writer) {
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.Usage()
}

func main() {

	var device string
	flag.StringVar(&device, "device", "", "device to read uart")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() {
		usage(os.Stderr)
		os.Exit(2)
	}

	run(device)
}
