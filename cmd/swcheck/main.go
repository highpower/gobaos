package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
	"strconv"
)

func run(device string, id string) error {
	value, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return err
	}
	conn, err := ft12.NewConn(device)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, baos.DatapointID(value))
}

func runWith(conn *ft12.Conn, id baos.DatapointID) error {
	req := baos.GetDatapointValueReq{StartDatapoint: id, NumberOfDatapoints: 1, Filter: baos.Filter(1)}
	data, err := req.MarshalBinary()
	if err != nil {
		return err
	}
	if _, err := conn.Write(data); err != nil {
		return err
	}
	for conn.Available() == 0 {
		if err := conn.Poll(); err != nil {
			return err
		}
	}
	buffer := make([]byte, conn.Available())
	if _, err := conn.Read(buffer); err != nil {
		return err
	}
	resp := baos.GetDatapointValueRes{}
	if err := resp.UnmarshalBinary(buffer); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", resp.String())
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: swcheck -device <device> <id>\n")
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.PrintDefaults()
}

func main() {

	var device string
	flag.StringVar(&device, "device", "", "device to connect to")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() || device == "" {
		usage(os.Stderr)
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) == 0 {
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := run(device, args[0]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
