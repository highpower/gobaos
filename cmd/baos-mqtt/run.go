package main

import (
	"errors"
	"fmt"
	"gobaos/internal/config"
	"gobaos/internal/mqttbridge"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func onIndication(atomicHandler atomic.Value, id baos.DatapointID, data []byte) error {
	if a := atomicHandler.Load(); a != nil {
		handler := a.(baos.ValueIndicationHandler)
		return handler.OnIndication(id, data)
	}
	return nil
}

func run(c *config.Data) error {
	log, err := newLogger(c)
	if err != nil {
		return fmt.Errorf("main.run: cannot create logger[%w]", err)
	}
	conn, err := ft12.NewConn(c.Baos.Device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return fmt.Errorf("main.run: cannot connect to baos device[%w]", err)
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, c, log)
}

func runWith(conn *ft12.Conn, c *config.Data, log *zap.SugaredLogger) error {
	atomicHandler := atomic.Value{}
	worker := baos.NewWorker(conn, baos.WithGoroutineCount(c.Baos.IndicationHandlers),
		baos.WithHandler(func(id baos.DatapointID, data []byte) error {
			return onIndication(atomicHandler, id, data)
		}))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go stopWorker(worker, signals)

	group := errgroup.Group{}
	group.Go(worker.Run)

	handler, err := mqttbridge.NewHandler(worker, c, mqttbridge.WithLogger(log))
	if err != nil {
		return errors.Join(fmt.Errorf("main.run: cannot create handler[%w]", err), group.Wait())
	}
	atomicHandler.Store(handler)
	log.Infof("main.run: baos-mqtt started with %s", c.String())
	return group.Wait()
}

func stopWorker(worker *baos.Worker, signals <-chan os.Signal) {
	<-signals
	worker.Stop()
}
