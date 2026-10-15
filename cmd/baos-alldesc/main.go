package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"

	"golang.org/x/sync/errgroup"
)

func run(device string) error {
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return fmt.Errorf("open device %q: %w", device, err)
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn)
}

func runWith(conn *ft12.Conn) error {
	group := errgroup.Group{}
	if err := runWithGroup(&group, conn); err != nil {
		return errors.Join(fmt.Errorf("runWithGroup: %w", err), group.Wait())
	}
	return group.Wait()
}

func runWithGroup(g *errgroup.Group, conn *ft12.Conn) error {
	worker := baos.NewWorker(conn)
	defer worker.Stop()
	g.Go(worker.Run)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	descriptions, err := worker.AllDescriptions(ctx)
	if err != nil {
		return fmt.Errorf("get all descriptions: %w", err)
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
