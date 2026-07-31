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

func run(device string, id string) error {
	value, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return err
	}
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, baos.DatapointID(value))
}

func runWorker(w *baos.Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := w.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func runWith(conn *ft12.Conn, id baos.DatapointID) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	worker := baos.NewWorker(conn)
	defer worker.Stop()

	wg.Add(1)
	go runWorker(worker, &wg)

	result, err := worker.Bool(id)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "%t\n", result)
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
