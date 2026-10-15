package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"

	"golang.org/x/sync/errgroup"
)

func run(device string, args []string) error {
	value, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return fmt.Errorf("parse datapoint id: %w", err)
	}
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return fmt.Errorf("open device %q: %w", device, err)
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, baos.DatapointID(value))
}

func runWith(conn *ft12.Conn, id baos.DatapointID) error {
	group := errgroup.Group{}
	if err := runWithGroup(&group, conn, id); err != nil {
		return errors.Join(fmt.Errorf("runWithGroup: %w", err), group.Wait())
	}
	return group.Wait()
}

func runWithGroup(g *errgroup.Group, conn *ft12.Conn, id baos.DatapointID) error {
	worker := baos.NewWorker(conn)
	defer worker.Stop()
	g.Go(worker.Run)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dpt, err := baos.CreateForID(ctx, worker, id)
	if err != nil {
		return fmt.Errorf("create dpt for datapoint %d: %w", id, err)
	}
	if err := worker.Get(ctx, id, dpt); err != nil {
		return fmt.Errorf("get datapoint %d: %w", id, err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: %s\n", id, dpt.Get())
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
	if len(args) != 1 {
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := run(device, args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
