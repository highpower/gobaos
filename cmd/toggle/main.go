package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"gobaos/pkg/knxtypes"
	"io"
	"os"
	"strconv"
	"strings"
)

func processRequest(conn *ft12.Conn, id baos.DatapointID, data []byte) error {
	req := baos.SetDatapointValueReq{Commands: []baos.DatapointCommand{{ID: id, Command: baos.SetNewValue, Data: data}}}
	reqData, err := req.MarshalBinary()
	if err != nil {
		return err
	}
	if _, err := conn.Write(reqData); err != nil {
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
	return processResponse(buffer)
}

func processResponse(buffer []byte) error {
	resp := baos.SetDatapointValueRes{}
	if err := resp.UnmarshalBinary(buffer); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", resp.String())
	return nil
}

func run(device string, args []string) error {
	id, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return err
	}
	on := true
	switch strings.ToLower(args[1]) {
	case "on":
		on = true
	case "off":
		on = false
	default:
		return fmt.Errorf("unknown command: %s", args[1])
	}
	return runWith(device, baos.DatapointID(id), on)
}

func runWith(device string, id baos.DatapointID, on bool) error {
	conn, err := ft12.NewConn(device)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	dpt := knxtypes.DPTBool{Value: on}
	data, err := dpt.MarshalBinary()
	if err != nil {
		return err
	}
	return processRequest(conn, id, data)
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: toggle -device <device> <id> <on|off>\n")
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
	if len(args) != 2 {
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := run(device, args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
