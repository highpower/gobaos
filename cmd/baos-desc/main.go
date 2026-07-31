package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

func run(device string, args []string) error {
	start, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return err
	}
	number, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {

	}
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, baos.DatapointID(start), uint16(number))
}

func runWorker(w *baos.Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := w.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func runWith(conn *ft12.Conn, start baos.DatapointID, number uint16) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	worker := baos.NewWorker(conn)
	defer worker.Stop()

	wg.Add(1)
	go runWorker(worker, &wg)

	descriptions, err := worker.Descriptions(start, number)
	if err != nil {
		return err
	}
	for _, d := range descriptions {
		_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: type %s\n", d.ID, d.DPT.String())
	}

	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: baos-get -device <device> <id>\n")
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
