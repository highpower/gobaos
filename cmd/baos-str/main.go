package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

func run(device string, args []string) error {
	id, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return err
	}
	conn, err := ft12.NewConn(device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, baos.DatapointID(id))
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

	s, err := worker.DescriptionString(ctx, id)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: type %s\n", id, s)
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: baos-str -device <device> <id>\n")
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
