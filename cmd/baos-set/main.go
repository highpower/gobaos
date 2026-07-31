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
	"sync"
	"time"
)

func run(device string, args []string) error {
	id, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return err
	}
	t, err := knxtypes.ParseTypeCode(args[1])
	if err != nil {
		return err
	}
	dpt, err := knxtypes.NewDPT(t)
	if err != nil {
		return err
	}
	if err := dpt.Set(args[2]); err != nil {
		return err
	}
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, baos.DatapointID(id), dpt)
}

func runWorker(w *baos.Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := w.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func runWith(conn *ft12.Conn, id baos.DatapointID, dpt knxtypes.DPT) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	worker := baos.NewWorker(conn)
	defer worker.Stop()

	wg.Add(1)
	go runWorker(worker, &wg)

	if err := worker.Set(id, dpt); err != nil {
		return err
	}
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: baos-set -device <device> <id> <on|off>\n")
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
	if len(args) != 3 {
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := run(device, args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
