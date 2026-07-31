package main

import (
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

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

func runWorker(w *baos.Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := w.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

}

func runWith(device string, id baos.DatapointID, on bool) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	conn, err := ft12.NewConn(device)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	worker := baos.NewWorker(conn)
	defer worker.Stop()
	wg.Add(1)
	go runWorker(worker, &wg)
	if err := worker.SetBool(id, on); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return err
	}
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: swtoggle -device <device> <id> <on|off>\n")
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
