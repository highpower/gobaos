package main

import (
	"context"
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
	"sync"
	"time"
)

func run(device string) error {
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn)
}

func runWorker(w *baos.Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := w.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func runWith(conn *ft12.Conn) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	worker := baos.NewWorker(conn)
	defer worker.Stop()

	wg.Add(1)
	go runWorker(worker, &wg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	descriptions, err := worker.AllDescriptions(ctx)
	if err != nil {
		return err
	}
	for _, d := range descriptions {
		_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: type %d-%s(%s)\n", d.ID, d.Code, d.Code.String(),
			baos.SizeString(d.ValueType))
	}
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: baos-desc -device <device> <id>\n")
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

	if err := run(device); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
