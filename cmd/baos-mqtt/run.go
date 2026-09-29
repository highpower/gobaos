package main

import (
	"fmt"
	"gobaos/internal/config"
	"gobaos/internal/mqttbridge"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"sync/atomic"
	"time"
)

func onIndication(pointer atomic.Value, id baos.DatapointID, data []byte) error {
	if a := pointer.Load(); a != nil {
		handler := a.(baos.ValueIndicationHandler)
		return handler.OnIndication(id, data)
	}
	return nil
}

func run(c *config.Data) error {
	logger, err := newLogger(c)
	if err != nil {
		return fmt.Errorf("main.run: cannot create logger[%w]", err)
	}
	conn, err := ft12.NewConn(c.Baos.Device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return fmt.Errorf("main.run: cannot connect to baos device[%w]", err)
	}
	defer func() { _ = conn.Close() }()
	value := atomic.Value{}
	worker := baos.NewWorker(conn,
		baos.WithGoroutineCount(c.Baos.IndicationHandlers),
		baos.WithHandler(func(id baos.DatapointID, data []byte) error {
			return onIndication(value, id, data)
		}))
	handler, err := mqttbridge.NewHandler(worker, c, mqttbridge.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("main.run: cannot create handler[%w]", err)
	}
	value.Store(handler)
	logger.Infof("main.run: baos-mqtt started with %s", c.String())
	return worker.Run()
}
